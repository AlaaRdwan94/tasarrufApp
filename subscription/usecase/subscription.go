package usecase

import (
	"context"
	"time"

	"github.com/ahmedaabouzied/tasarruf/branch"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/offer"
	"github.com/ahmedaabouzied/tasarruf/payment"
	"github.com/ahmedaabouzied/tasarruf/subscription"
	"github.com/ahmedaabouzied/tasarruf/user"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// SubscriptionUsecase implements subscription usecase interface
type SubscriptionUsecase struct {
	SubscriptionRepo subscription.Repository
	UserRepo         user.Repository
	BranchRepo       branch.Repository
	OfferRepo        offer.Repository
}

// CreateSubscriptionUsecase returns an implementation of the subscription usecase interface
func CreateSubscriptionUsecase(subscriptionRepo subscription.Repository, userRepo user.Repository, branchRepo branch.Repository, offerRepo offer.Repository) subscription.Usecase {
	u := SubscriptionUsecase{
		SubscriptionRepo: subscriptionRepo,
		UserRepo:         userRepo,
		BranchRepo:       branchRepo,
		OfferRepo:        offerRepo,
	}
	return &u
}

// CreatePlan creates a new subscription plan
func (u *SubscriptionUsecase) CreatePlan(ctx context.Context, p *entities.Plan) (*entities.Plan, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user")
		log.Error(err, "current user id = ", currentUserID)
		cancelFunc()
		return nil, err
	}
	if currentUser.AccountType != "admin" {
		err := errors.New("only admin users can create subscription plans")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if p.IsDefault {
		defaultPlan, err := u.SubscriptionRepo.GetDefaultPlan(ctx)
		if err != nil {
			if errors.Cause(err).Error() == "plan not found" {
				log.Error(err)
			} else {
				cancelFunc()
				return nil, err
			}
		} else {
			defaultPlan.IsDefault = false
			_, err = u.SubscriptionRepo.UpdatePlan(ctx, defaultPlan)
			if err != nil {
				cancelFunc()
				return nil, err
			}
		}
	}
	createdPlan, err := u.SubscriptionRepo.CreatePlan(ctx, p)
	if err != nil {
		err = errors.Wrap(err, "repository error while creating plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return createdPlan, nil
}

// DeletePlan deletes the plan with the given ID
func (u *SubscriptionUsecase) DeletePlan(ctx context.Context, planID uint) (*entities.Plan, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user")
		log.Error(err, "current user id = ", currentUserID)
		cancelFunc()
		return nil, err
	}
	if currentUser.AccountType != "admin" {
		err := errors.New("only admin users can delete subscription plans")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	toDeletePlan, err := u.SubscriptionRepo.GetPlanByID(ctx, planID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if toDeletePlan.IsDefault {
		err := errors.New("cannot delete the default plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	deletedPlan, err := u.SubscriptionRepo.DeletePlan(ctx, toDeletePlan)
	if err != nil {
		err = errors.Wrap(err, "repository error while deleting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return deletedPlan, nil
}

// UpdatePlan updates the plan with the given ID
func (u *SubscriptionUsecase) UpdatePlan(ctx context.Context, planID uint, plan *entities.Plan) (*entities.Plan, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user")
		log.Error(err, "current user id = ", currentUserID)
		cancelFunc()
		return nil, err
	}
	if currentUser.AccountType != "admin" {
		err := errors.New("only admin users can delete subscription plans")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	toUpdatePlan, err := u.SubscriptionRepo.GetPlanByID(ctx, planID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	toUpdatePlan.ID = planID
	toUpdatePlan.EnglishName = plan.EnglishName
	toUpdatePlan.TurkishName = plan.TurkishName
	toUpdatePlan.EnglishDescription = plan.EnglishDescription
	toUpdatePlan.TurkishDescription = plan.TurkishDescription
	toUpdatePlan.CountOfOffers = plan.CountOfOffers
	toUpdatePlan.Price = plan.Price
	toUpdatePlan.Image = plan.Image
	if toUpdatePlan.IsDefault != plan.IsDefault {
		defaultPlan, err := u.SubscriptionRepo.GetDefaultPlan(ctx)
		if err != nil {
			if errors.Cause(err).Error() == "plan not found" {
				log.Error(err)
			} else {
				cancelFunc()
				return nil, err
			}
		} else {
			defaultPlan.IsDefault = false
			_, err = u.SubscriptionRepo.UpdatePlan(ctx, defaultPlan)
			if err != nil {
				cancelFunc()
				return nil, err
			}
		}
		toUpdatePlan.IsDefault = plan.IsDefault
	}
	updatedPlan, err := u.SubscriptionRepo.UpdatePlan(ctx, toUpdatePlan)
	if err != nil {
		err = errors.Wrap(err, "repository error while deleting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return updatedPlan, nil

}

// RankPlanUp ranks the plan with the given ID up
func (u *SubscriptionUsecase) RankPlanUp(ctx context.Context, planID uint) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	plan, err := u.SubscriptionRepo.GetPlanByID(ctx, planID)
	if err != nil {
		cancelFunc()
		return err
	}
	plan.Rank = plan.Rank + 1
	_, err = u.SubscriptionRepo.RankPlanUp(ctx, plan)
	if err != nil {
		cancelFunc()
		return err
	}
	cancelFunc()
	return nil
}

// GetAllPlans returns all the subscription plans
func (u *SubscriptionUsecase) GetAllPlans(ctx context.Context) ([]entities.Plan, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting current user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	plans, err := u.SubscriptionRepo.GetPlans(ctx)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plans")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if !currentUser.IsAdmin() {
		var filtered []entities.Plan
		for _, plan := range plans {
			if !plan.IsDefault {
				filtered = append(filtered, plan)
			}
		}
		cancelFunc()
		return filtered, nil
	}
	cancelFunc()
	return plans, err
}

// SubscribeToPlan subscribes the current user to the plan with the given ID
func (u *SubscriptionUsecase) SubscribeToPlan(ctx context.Context, planID uint, paymentDetails *entities.PaymentDetails) (*entities.Subscription, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	userID := ctx.Value(entities.UserIDKey).(uint)
	user, err := u.UserRepo.GetByID(ctx, userID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if user.AccountType == "partner" {
		err = errors.New("partner users cannot subscribe to plans")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	userCity, err := u.BranchRepo.GetCityByID(ctx, user.CityID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	user.City = *userCity
	userCurrentPlan, err := u.SubscriptionRepo.GetSubscriptionByUser(ctx, user.ID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user current plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if userCurrentPlan != nil {
		err = errors.New("user is already registered to a plan")
		log.Error(err)
		cancelFunc()
		return nil, err

	}
	plan, err := u.SubscriptionRepo.GetPlanByID(ctx, planID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err

	}
	paymentDetails.User = user
	paymentDetails.Plan = plan
	p := payment.CreateTransaction(paymentDetails)
	id, err := p.Submit(ctx)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	log.Info("payment id = ", id)
	subscription := &entities.Subscription{
		UserID:              user.ID,
		PlanID:              plan.ID,
		RemainingOffers:     plan.CountOfOffers,
		Expired:             false,
		ExpireDate:          time.Now().AddDate(1, 0, 0),
		DelegationStartDate: time.Now(),
	}
	subscription, err = u.SubscriptionRepo.CreateSubscription(ctx, subscription)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	subscription.Plan = *plan
	cancelFunc()
	return subscription, nil
}

// UpgradePlan subscribes the user to a new plan while the remaining offers get added to  new subscription.
func (u *SubscriptionUsecase) UpgradePlan(ctx context.Context, planID uint, paymentDetails *entities.PaymentDetails) (*entities.Subscription, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	userID := ctx.Value(entities.UserIDKey).(uint)
	customer, err := u.getCustomerByID(ctx, userID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while getting customer")
	}
	userCity, err := u.BranchRepo.GetCityByID(ctx, customer.CityID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	customer.City = *userCity
	err = customer.ExpireSubscription()
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	newPlan, err := u.SubscriptionRepo.GetPlanByID(ctx, planID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	paymentDetails.User = &customer.User
	paymentDetails.Plan = newPlan
	p := payment.CreateTransaction(paymentDetails)
	id, err := p.Submit(ctx)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if id == "" {
		err = errors.Wrap(err, "error processing payment")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	oldSubscription := customer.Subscription
	oldCountsOfOffers, err := u.SubscriptionRepo.GetCountOfOffersOfCustomer(ctx, oldSubscription)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}

	subscription := &entities.Subscription{
		UserID:              customer.ID,
		PlanID:              newPlan.ID,
		RemainingOffers:     customer.Subscription.GetRemainingOffers() + newPlan.CountOfOffers,
		Expired:             false,
		ExpireDate:          time.Now().AddDate(1, 0, 0),
		DelegationStartDate: customer.Subscription.DelegationStartDate,
	}
	subscription, err = u.SubscriptionRepo.CreateSubscription(ctx, subscription)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	subscription.Plan = *newPlan
	_, err = u.SubscriptionRepo.ExpireSubscription(ctx, customer.Subscription)
	if err != nil {
		err = errors.Wrap(err, "repository error while upgrading subscription")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	for _, countOfOffer := range oldCountsOfOffers {
		partner, err := u.getPartnerByID(ctx, countOfOffer.PartnerID)
		if err != nil {
			log.Error(err)
			cancelFunc()
		}
		// set new counts of offers
		err = u.SubscriptionRepo.SetCountOfOffersWithPartner(ctx, partner, subscription, countOfOffer.CountOfOffers+newPlan.CountOfOffers)
		if err != nil {
			log.Error(err)
		}
	}
	cancelFunc()
	return subscription, nil
}

// RenewPlan creates a new subscription for the user with the plan currently subscribed to.
// While it sets the old subscription as expired.
// Sets the remaining offers of the new plan to the sum of remaining offers of the old plan and the offers of the new plan.
func (u *SubscriptionUsecase) RenewPlan(ctx context.Context, paymentDetails *entities.PaymentDetails) (*entities.Subscription, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	userID := ctx.Value(entities.UserIDKey).(uint)
	user, err := u.UserRepo.GetByID(ctx, userID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if user.AccountType == "partner" {
		err = errors.New("partner users cannot subscribe to plans")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	userCity, err := u.BranchRepo.GetCityByID(ctx, user.CityID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	user.City = *userCity

	userCurrentSubscription, err := u.SubscriptionRepo.GetSubscriptionByUser(ctx, user.ID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user current subscription")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	oldSubscription := userCurrentSubscription
	oldCountsOfOffers, err := u.SubscriptionRepo.GetCountOfOffersOfCustomer(ctx, oldSubscription)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if userCurrentSubscription == nil {
		err = errors.New("user is not subscribed to any plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	plan, err := u.SubscriptionRepo.GetPlanByID(ctx, userCurrentSubscription.PlanID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	paymentDetails.User = user
	paymentDetails.Plan = plan
	p := payment.CreateTransaction(paymentDetails)
	_, err = p.Submit(ctx)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	subscription := &entities.Subscription{
		UserID:              user.ID,
		PlanID:              plan.ID,
		RemainingOffers:     plan.CountOfOffers + userCurrentSubscription.RemainingOffers,
		Expired:             false,
		ExpireDate:          time.Now().AddDate(1, 0, 0),
		DelegationStartDate: userCurrentSubscription.DelegationStartDate,
	}
	subscription, err = u.SubscriptionRepo.CreateSubscription(ctx, subscription)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	subscription.Plan = *plan
	_, err = u.SubscriptionRepo.ExpireSubscription(ctx, userCurrentSubscription)
	if err != nil {
		err = errors.Wrap(err, "repository error while upgrading subscription")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	for _, countOfOffer := range oldCountsOfOffers {
		partner, err := u.getPartnerByID(ctx, countOfOffer.PartnerID)
		if err != nil {
			log.Error(err)
			cancelFunc()
		}
		// set new counts of offers
		err = u.SubscriptionRepo.SetCountOfOffersWithPartner(ctx, partner, subscription, countOfOffer.CountOfOffers+plan.CountOfOffers)
		if err != nil {
			log.Error(err)
		}
	}
	cancelFunc()
	return subscription, nil
}

// GetMySubscription returns the current subscription of the currently logged in user
func (u *SubscriptionUsecase) GetMySubscription(ctx context.Context) (*entities.Subscription, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	userID := ctx.Value(entities.UserIDKey).(uint)
	user, err := u.UserRepo.GetByID(ctx, userID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if user.AccountType == "partner" {
		err = errors.New("partner users cannot subscribe to plans")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	userCurrentSubscription, err := u.SubscriptionRepo.GetSubscriptionByUser(ctx, user.ID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user current subscription")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if userCurrentSubscription == nil {
		err = errors.New("user is not subscribed to any plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	plan, err := u.SubscriptionRepo.GetPlanByID(ctx, userCurrentSubscription.PlanID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user current subscription plan")
		log.Error(err)
		cancelFunc()
		return nil, err

	}
	userCurrentSubscription.Plan = *plan
	cancelFunc()
	return userCurrentSubscription, nil
}

// AdminUpgradeUserPlan upgared user plan by admin
func (u *SubscriptionUsecase) AdminUpgradeUserPlan(ctx context.Context, userID uint, planID uint) (*entities.Subscription, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting current user")
	}
	if !currentUser.IsAdmin() {
		cancelFunc()
		return nil, errors.New("not authorized to change user's plan")
	}
	customer, err := u.getCustomerByID(ctx, userID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while getting customer")
	}
	userCity, err := u.BranchRepo.GetCityByID(ctx, customer.CityID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	customer.City = *userCity
	if customer.Subscription != nil {
		err = customer.ExpireSubscription()
		if err != nil {
			log.Error(err)
			cancelFunc()
			return nil, err
		}
	}
	oldSubscription := customer.Subscription
	oldCountsOfOffers, err := u.SubscriptionRepo.GetCountOfOffersOfCustomer(ctx, oldSubscription)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	newPlan, err := u.SubscriptionRepo.GetPlanByID(ctx, planID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	subscription := &entities.Subscription{
		UserID:     customer.ID,
		PlanID:     newPlan.ID,
		Expired:    false,
		ExpireDate: time.Now().AddDate(1, 0, 0),
	}
	if customer.Subscription != nil {
		subscription.RemainingOffers = customer.Subscription.GetRemainingOffers() + newPlan.CountOfOffers
		subscription.DelegationStartDate = customer.Subscription.DelegationStartDate
	} else {
		subscription.RemainingOffers = newPlan.CountOfOffers
		subscription.DelegationStartDate = time.Now()
	}
	subscription, err = u.SubscriptionRepo.CreateSubscription(ctx, subscription)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	subscription.Plan = *newPlan
	if customer.Subscription != nil {
		_, err = u.SubscriptionRepo.ExpireSubscription(ctx, customer.Subscription)
		if err != nil {
			err = errors.Wrap(err, "repository error while upgrading subscription")
			log.Error(err)
			cancelFunc()
			return nil, err
		}
	}
	for _, countOfOffer := range oldCountsOfOffers {
		partner, err := u.getPartnerByID(ctx, countOfOffer.PartnerID)
		if err != nil {
			log.Error(err)
			cancelFunc()
		}
		// set new counts of offers
		err = u.SubscriptionRepo.SetCountOfOffersWithPartner(ctx, partner, subscription, countOfOffer.CountOfOffers+newPlan.CountOfOffers)
		if err != nil {
			log.Error(err)
		}
	}
	cancelFunc()
	return subscription, nil
}

func (u *SubscriptionUsecase) SubscribeToFreePlan(ctx context.Context, planID uint) (*entities.Subscription, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting current user")
	}
	customer, err := u.getCustomerByID(ctx, currentUser.GetID())
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while getting customer")
	}
	userCity, err := u.BranchRepo.GetCityByID(ctx, customer.CityID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	customer.City = *userCity
	if customer.Subscription != nil {
		err = customer.ExpireSubscription()
		if err != nil {
			log.Error(err)
			cancelFunc()
			return nil, err
		}
	}
	newPlan, err := u.SubscriptionRepo.GetPlanByID(ctx, planID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	subscription := &entities.Subscription{
		UserID:     customer.ID,
		PlanID:     newPlan.ID,
		Expired:    false,
		ExpireDate: time.Now().AddDate(1, 0, 0),
	}
	if customer.Subscription != nil {
		subscription.RemainingOffers = customer.Subscription.GetRemainingOffers() + newPlan.CountOfOffers
		subscription.DelegationStartDate = customer.Subscription.DelegationStartDate
	} else {
		subscription.RemainingOffers = newPlan.CountOfOffers
		subscription.DelegationStartDate = time.Now()
	}
	subscription, err = u.SubscriptionRepo.CreateSubscription(ctx, subscription)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	subscription.Plan = *newPlan
	if customer.Subscription != nil {
		_, err = u.SubscriptionRepo.ExpireSubscription(ctx, customer.Subscription)
		if err != nil {
			err = errors.Wrap(err, "repository error while upgrading subscription")
			log.Error(err)
			cancelFunc()
			return nil, err
		}
	}
	cancelFunc()
	return subscription, nil
}

// GetMySubscriptionWithPartner returns the subscription with the count of offers for the given partner
func (u *SubscriptionUsecase) GetMySubscriptionWithPartner(ctx context.Context, partnerID uint) (*entities.Subscription, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	userID := ctx.Value(entities.UserIDKey).(uint)
	user, err := u.UserRepo.GetByID(ctx, userID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if user.AccountType == "partner" {
		err = errors.New("partner users cannot subscribe to plans")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	userCurrentSubscription, err := u.SubscriptionRepo.GetSubscriptionByUser(ctx, user.ID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user current subscription")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if userCurrentSubscription == nil {
		err = errors.New("user is not subscribed to any plan")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	plan, err := u.SubscriptionRepo.GetPlanByID(ctx, userCurrentSubscription.PlanID)
	if err != nil {
		err = errors.Wrap(err, "repository error while getting user current subscription")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	userCurrentSubscription.Plan = *plan
	partner, err := u.getPartnerByID(ctx, partnerID)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	cof, err := u.SubscriptionRepo.GetCountOfOffersWithPartner(ctx, partner, userCurrentSubscription)
	if err != nil {
		err = errors.New("error getting user remaining offers")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	userCurrentSubscription.RemainingOffers = cof.CountOfOffers
	userCurrentSubscription.Plan = *plan
	cancelFunc()
	return userCurrentSubscription, nil

}

// CreatePlanCategoryAssociation creates a plan-category association
func (u *SubscriptionUsecase) CreatePlanCategoryAssociation(ctx context.Context, planID uint, categoryID uint) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		cancelFunc()
		return errors.New("not authorized")
	}
	err = u.SubscriptionRepo.CreatePlanCategoryAssociation(ctx, planID, categoryID)
	if err != nil {
		cancelFunc()
		return err
	}
	cancelFunc()
	return nil
}

// RemovePlanCategoryAssociation removes the plan-category association
func (u *SubscriptionUsecase) RemovePlanCategoryAssociation(ctx context.Context, planID uint, categoryID uint) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		cancelFunc()
		return errors.New("not authorized")
	}
	err = u.SubscriptionRepo.RemovePlanCategoryAssociation(ctx, planID, categoryID)
	if err != nil {
		cancelFunc()
		return err
	}
	cancelFunc()
	return nil
}

// GetCategoriesOfPlan returns the categories associated for the given plan
func (u *SubscriptionUsecase) GetCategoriesOfPlan(ctx context.Context, planID uint) ([]entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	categories, err := u.SubscriptionRepo.GetCategoriesByPlanID(ctx, planID)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return categories, nil
}

// GetPlanByID returns the plan with the given ID. It returns an error if the plan is not found.
func (u *SubscriptionUsecase) GetPlanByID(ctx context.Context, planID uint) (*entities.Plan, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	plan, err := u.SubscriptionRepo.GetPlanByID(ctx, planID)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return plan, nil

}

func (u *SubscriptionUsecase) getCustomerByID(ctx context.Context, ID uint) (*entities.Customer, error) {
	user, err := u.UserRepo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	otp, err := u.UserRepo.GetOTPByUser(ctx, ID)
	if err != nil {
		return nil, err
	}
	subscription, err := u.SubscriptionRepo.GetSubscriptionByUser(ctx, ID)
	if err != nil {
		return nil, err
	}
	if subscription != nil {
		plan, err := u.SubscriptionRepo.GetPlanByID(ctx, subscription.PlanID)
		if err != nil {
			return nil, err
		}
		subscription.Plan = *plan
	}
	customer := entities.Customer{
		User:         *user,
		OTP:          otp,
		Subscription: subscription,
	}
	return &customer, nil
}

func (u *SubscriptionUsecase) getPartnerByID(ctx context.Context, ID uint) (*entities.Partner, error) {
	user, err := u.UserRepo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	otp, err := u.UserRepo.GetOTPByUser(ctx, ID)
	if err != nil {
		return nil, err
	}
	partnerProfile, err := u.UserRepo.GetPartnerProfileByUserID(ctx, ID)
	if err != nil {
		return nil, err
	}
	partner := entities.Partner{
		User:           *user,
		OTP:            otp,
		PartnerProfile: partnerProfile,
	}
	return &partner, nil
}
