package review

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
)

// Repository defines the contract of the review repository
type Repository interface {
	Create(ctx context.Context, review *entities.Review) (*entities.Review, error)
	Update(ctx context.Context, review *entities.Review) (*entities.Review, error)
	Delete(ctx context.Context, review *entities.Review) (*entities.Review, error)
	GetByID(ctx context.Context, ID uint) (*entities.Review, error)
	GetByCustomerID(ctx context.Context, ID uint) ([]entities.Review, error)
	GetByPartnerID(ctx context.Context, ID uint) ([]entities.Review, error)
	GetAverageRatings(ctx context.Context, partnerID uint) (float64, error)
}
