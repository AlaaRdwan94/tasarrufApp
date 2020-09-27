package subscription

import (
	"context"

	"github.com/ahmedaabouzied/tasarruf/entities"
)

// Repository represents the contract of the user repository
type Repository interface {
	CreatePlan(ctx context.Context, p *entities.Plan) (*entities.Plan, error)
	GetPlanByID(ctx context.Context, planID uint) (*entities.Plan, error)
	UpdatePlan(ctx context.Context, plan *entities.Plan) (*entities.Plan, error)
	DeletePlan(ctx context.Context, p *entities.Plan) (*entities.Plan, error)
	GetPlans(ctx context.Context) ([]entities.Plan, error)
	CreateSubscription(ctx context.Context, s *entities.Subscription) (*entities.Subscription, error)
	DeleteSubscription(ctx context.Context, s *entities.Subscription) (*entities.Subscription, error)
	DecrementOffersCount(ctx context.Context, s *entities.Subscription) (*entities.Subscription, error)
	GetSubscriptionByID(ctx context.Context, ID uint) (*entities.Subscription, error)
	GetSubscriptionByUser(ctx context.Context, userID uint) (*entities.Subscription, error)
	ExpireSubscription(ctx context.Context, s *entities.Subscription) (*entities.Subscription, error)
	RankPlanUp(ctx context.Context, toUpdatePlan *entities.Plan) ([]entities.Plan, error)
	CreatePlanCategoryAssociation(ctx context.Context, planID uint, CategoryID uint) error
	GetCategoriesByPlanID(ctx context.Context, planID uint) ([]entities.Category, error)
	RemovePlanCategoryAssociation(ctx context.Context, planID uint, CategoryID uint) error
	GetDefaultPlan(ctx context.Context) (*entities.Plan, error)
	GetCountOfOffersWithPartner(ctx context.Context, partner *entities.Partner, subscription *entities.Subscription) (*entities.CustomerPartnerOffersCount, error)
	SetCountOfOffersWithPartner(ctx context.Context, partner *entities.Partner, subscription *entities.Subscription, newCount uint) error
	GetCountOfOffersOfCustomer(ctx context.Context, subscription *entities.Subscription) ([]entities.CustomerPartnerOffersCount, error)
}
