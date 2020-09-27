package repository

import (
	"context"
	"time"

	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/subscription"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// SubscriptionRepository implements subscription repository interface
type SubscriptionRepository struct {
	DB *gorm.DB
}

// CreateSubscriptionRepository return an implementation of subscription repository interface
func CreateSubscriptionRepository(db *gorm.DB) subscription.Repository {
	repo := SubscriptionRepository{
		DB: db,
	}
	return &repo
}

// CreatePlan creates a new subscription plan
func (r *SubscriptionRepository) CreatePlan(ctx context.Context, p *entities.Plan) (*entities.Plan, error) {
	dbt := r.DB.Create(p)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error creating plan")
	}
	return p, nil
}

// GetPlanByID returns the plan with the given ID
func (r *SubscriptionRepository) GetPlanByID(ctx context.Context, planID uint) (*entities.Plan, error) {
	var plan entities.Plan
	dbt := r.DB.Where("id = ?", planID).Find(&plan)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return nil, errors.New("plan not found")
		}
		return nil, errors.Wrap(dbt.Error, "error getting plan")
	}
	return &plan, nil
}

// GetDefaultPlan returns the plan with the default property
func (r *SubscriptionRepository) GetDefaultPlan(ctx context.Context) (*entities.Plan, error) {
	var plan entities.Plan
	dbt := r.DB.Where("is_default = true").Find(&plan)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return nil, errors.New("plan not found")
		}
		return nil, errors.Wrap(dbt.Error, "error getting plan")
	}
	return &plan, nil
}

// DeletePlan deletes the given subscription plan
func (r *SubscriptionRepository) DeletePlan(ctx context.Context, p *entities.Plan) (*entities.Plan, error) {
	dbt := r.DB.Delete(p)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error deleting plan")
	}
	return p, nil
}

// UpdatePlan saves the given plan
func (r *SubscriptionRepository) UpdatePlan(ctx context.Context, plan *entities.Plan) (*entities.Plan, error) {
	dbt := r.DB.Save(plan)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error updating plan")
	}
	return plan, nil
}

// RankPlanUp sets the ranking to the given plan up by one level
func (r *SubscriptionRepository) RankPlanUp(ctx context.Context, toUpdatePlan *entities.Plan) ([]entities.Plan, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var toUpdatePlans []entities.Plan
	allPlans, err := r.GetPlans(ctx)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	for _, plan := range allPlans {
		if plan.Rank == toUpdatePlan.Rank-1 {
			plan.Rank = toUpdatePlan.Rank
			toUpdatePlans = append(toUpdatePlans, plan)
		} else {
			if plan.Rank == toUpdatePlan.Rank {
				plan.Rank = toUpdatePlan.Rank - 1
				toUpdatePlans = append(toUpdatePlans, plan)
			}
		}
	}
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		cancelFunc()
		return nil, err
	}

	for _, plan := range toUpdatePlans {
		tx.Save(plan)
	}

	cancelFunc()
	return toUpdatePlans, tx.Commit().Error
}

// GetPlans returns all plan records
func (r *SubscriptionRepository) GetPlans(ctx context.Context) ([]entities.Plan, error) {
	var plans []entities.Plan
	dbt := r.DB.Order("rank ASC").Find(&plans)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.Wrap(dbt.Error, "error getting plans")
	}
	return plans, nil
}

// CreateSubscription creates a new subscription record
func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, s *entities.Subscription) (*entities.Subscription, error) {
	dbt := r.DB.Create(s)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error creating subscription record")
	}
	return s, nil
}

// DeleteSubscription soft deletes the subscription record
func (r *SubscriptionRepository) DeleteSubscription(ctx context.Context, s *entities.Subscription) (*entities.Subscription, error) {
	dbt := r.DB.Delete(s)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error creating subscription record")
	}
	return s, nil
}

// DecrementOffersCount decrements the count of remaining offers in the given subscription
func (r *SubscriptionRepository) DecrementOffersCount(ctx context.Context, s *entities.Subscription) (*entities.Subscription, error) {
	dbt := r.DB.Save(s)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error saving change to subscription")
	}
	return s, nil
}

// GetSubscriptionByID returns the subscription with the given ID
func (r *SubscriptionRepository) GetSubscriptionByID(ctx context.Context, ID uint) (*entities.Subscription, error) {
	var subscription entities.Subscription
	dbt := r.DB.Where("id = ?", ID).Find(&subscription)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.Wrap(dbt.Error, "error getting subscription from DB")
	}
	return &subscription, nil
}

// GetSubscriptionByUser return the active subscription of the given user
func (r *SubscriptionRepository) GetSubscriptionByUser(ctx context.Context, userID uint) (*entities.Subscription, error) {
	var subscription entities.Subscription
	dbt := r.DB.Where("user_id = ? AND expired = ?", userID, false).Find(&subscription)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			defaultPlan, err := r.GetDefaultPlan(ctx)
			if err != nil {
				return nil, err
			}
			newSubscription := &entities.Subscription{
				UserID:              userID,
				PlanID:              defaultPlan.ID,
				RemainingOffers:     defaultPlan.CountOfOffers,
				PaymentID:           "",
				Expired:             false,
				ExpireDate:          time.Now().AddDate(1, 0, 0),
				DelegationStartDate: time.Now(),
			}
			newSubscription, err = r.CreateSubscription(ctx, newSubscription)
			if err != nil {
				return nil, err
			}
			return newSubscription, nil
		}
		return nil, errors.Wrap(dbt.Error, "error getting subscription of the given user")
	}
	if subscription.HasExpirPassed() {
		newSubscription, err := r.ExpireSubscription(ctx, &subscription)
		if err != nil {
			return nil, err
		}
		return newSubscription, nil
	}
	return &subscription, nil
}

// ExpireSubscription sets the given subscription as expired
func (r *SubscriptionRepository) ExpireSubscription(ctx context.Context, s *entities.Subscription) (*entities.Subscription, error) {
	s.Expired = true
	dbt := r.DB.Save(s)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error saving change to subscription")
	}
	return s, nil
}

// CreatePlanCategoryAssociation creates a plan-category association
func (r *SubscriptionRepository) CreatePlanCategoryAssociation(ctx context.Context, planID uint, CategoryID uint) error {
	pc := entities.PlanCategory{
		PlanID:     planID,
		CategoryID: CategoryID,
	}
	dbt := r.DB.Create(&pc)
	if dbt.Error != nil {
		return errors.Wrap(dbt.Error, "error creating plan category association")
	}
	return nil
}

// GetCategoriesByPlanID returns the categories associated with the given plan
func (r *SubscriptionRepository) GetCategoriesByPlanID(ctx context.Context, planID uint) ([]entities.Category, error) {
	var categories []entities.Category
	dbt := r.DB.Raw(`
        SELECT * FROM categories WHERE id IN (
            SELECT category_id FROM plan_categories WHERE plan_id = ?
        )
    `, planID).Scan(&categories)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting categories of the given plan")
	}
	return categories, nil
}

// RemovePlanCategoryAssociation removes the association of plan and category
func (r *SubscriptionRepository) RemovePlanCategoryAssociation(ctx context.Context, planID uint, CategoryID uint) error {
	dbt := r.DB.Exec(`DELETE FROM plan_categories WHERE plan_id = ? AND category_id = ?`, planID, CategoryID)
	if dbt.Error != nil {
		return errors.Wrap(dbt.Error, "error removing plan category association")
	}
	return nil
}

func (r *SubscriptionRepository) SetCountOfOffersWithPartner(ctx context.Context, partner *entities.Partner, subscription *entities.Subscription, newCount uint) error {
	cpoc, err := r.GetCountOfOffersWithPartner(ctx, partner, subscription)
	if err != nil {
		return err
	}
	cpoc.CountOfOffers = newCount
	dbt := r.DB.Save(cpoc)
	if dbt.Error != nil {
		return dbt.Error
	}
	return nil
}

func (r *SubscriptionRepository) GetCountOfOffersWithPartner(ctx context.Context, partner *entities.Partner, subscription *entities.Subscription) (*entities.CustomerPartnerOffersCount, error) {
	var customerPartnerOffersRecord entities.CustomerPartnerOffersCount
	dbt := r.DB.Where("customer_id = ? AND partner_id = ? AND subscription_id = ?", subscription.UserID, partner.ID, subscription.ID).Find(&customerPartnerOffersRecord)
	if dbt.RecordNotFound() {
		cpoc := &entities.CustomerPartnerOffersCount{
			CustomerID:     subscription.UserID,
			PartnerID:      partner.ID,
			SubscriptionID: subscription.ID,
		}
		if partner.PartnerProfile.IsSharable {
			var s entities.Share
			dbt := r.DB.Where("customer_id = ?", subscription.UserID).First(&s)
			if dbt.RecordNotFound() {
				cpoc.CountOfOffers = subscription.RemainingOffers
			} else {
				if s.ID > 0 {
					cpoc.CountOfOffers = subscription.RemainingOffers + 1
				} else {
					cpoc.CountOfOffers = subscription.RemainingOffers
				}
			}
		} else {
			cpoc.CountOfOffers = subscription.RemainingOffers
		}
		dbt := r.DB.Create(cpoc)
		if dbt.Error != nil {
			return cpoc, dbt.Error
		}
		return cpoc, nil
	}
	return &customerPartnerOffersRecord, nil
}

// GetCountOfOffersOfCustomer returns a list of couts of offers of the given customer
func (r *SubscriptionRepository) GetCountOfOffersOfCustomer(ctx context.Context, subscription *entities.Subscription) ([]entities.CustomerPartnerOffersCount, error) {
	var customerPartnerOffersRecords []entities.CustomerPartnerOffersCount
	dbt := r.DB.Where("customer_id = ? AND subscription_id = ?", subscription.UserID, subscription.ID).Find(&customerPartnerOffersRecords)
	if dbt.RecordNotFound() {
		return nil, nil
	}
	if dbt.Error != nil {
		return nil, dbt.Error
	}
	return customerPartnerOffersRecords, nil
}
