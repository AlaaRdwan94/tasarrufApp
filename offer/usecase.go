package offer

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"time"
)

// Usecase reporesents offer usecase contract
type Usecase interface {
	ConsumeOffer(ctx context.Context, customerID uint, partnerID uint, amount float64) (*entities.Offer, error)
	GetMyOffersHistory(ctx context.Context, startDate time.Time, endDate time.Time) ([]entities.Offer, error)
	SendOffersStaticMail(ctx context.Context, startDate time.Time, endDate time.Time) error
	GetOffer(ctx context.Context, offerID uint) (*entities.Offer, error)
	GetOffersCount(ctx context.Context) (int, error)
	GetAllOffers(ctx context.Context) ([]entities.Offer, error)
	GetByCustomer(ctx context.Context, customerID uint) ([]entities.Offer, error)
}
