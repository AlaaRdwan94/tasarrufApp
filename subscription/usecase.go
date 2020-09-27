package subscription

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
)

// Usecase implements the usecase contract of subscriptions
type Usecase interface {
	CreatePlan(ctx context.Context, p *entities.Plan) (*entities.Plan, error)
	DeletePlan(ctx context.Context, planID uint) (*entities.Plan, error)
	UpdatePlan(ctx context.Context, planID uint, plan *entities.Plan) (*entities.Plan, error)
	GetAllPlans(ctx context.Context) ([]entities.Plan, error)
	GetPlanByID(ctx context.Context, id uint) (*entities.Plan, error)
	SubscribeToPlan(ctx context.Context, planID uint, paymentDetails *entities.PaymentDetails) (*entities.Subscription, error)
	UpgradePlan(ctx context.Context, planID uint, paymentDetails *entities.PaymentDetails) (*entities.Subscription, error)
	RenewPlan(ctx context.Context, paymentDetails *entities.PaymentDetails) (*entities.Subscription, error)
	GetMySubscription(ctx context.Context) (*entities.Subscription, error)
	GetMySubscriptionWithPartner(ctx context.Context, partnerID uint) (*entities.Subscription, error)
	RankPlanUp(ctx context.Context, planID uint) error
	AdminUpgradeUserPlan(ctx context.Context, userID uint, planID uint) (*entities.Subscription, error)
	CreatePlanCategoryAssociation(ctx context.Context, planID uint, categoryID uint) error
	RemovePlanCategoryAssociation(ctx context.Context, planID uint, categoryID uint) error
	GetCategoriesOfPlan(ctx context.Context, planID uint) ([]entities.Category, error)
	SubscribeToFreePlan(ctx context.Context, planID uint) (*entities.Subscription, error)
}
