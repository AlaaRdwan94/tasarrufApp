// Copyright 2019 NOVA Solutions Co. All Rights Reserved.

package offer

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"time"
)

// Repository defines the contract of the offer repository
type Repository interface {
	Create(ctx context.Context, u *entities.Offer) (*entities.Offer, error)
	GetByID(ctx context.Context, ID uint) (*entities.Offer, error)
	GetByUser(ctx context.Context, userID uint) ([]entities.Offer, error)
	GetByPartner(ctx context.Context, partnerID uint, startDate time.Time, endDate time.Time) ([]entities.Offer, error)
	GetCountByPartnerAndCustomer(ctx context.Context, partnerID uint, customerID uint, startDate time.Time, endDate time.Time) (int, error)
	GetOffersCount(ctx context.Context) (int, error)
	GetAllOffers(ctx context.Context) ([]entities.Offer, error)
}
