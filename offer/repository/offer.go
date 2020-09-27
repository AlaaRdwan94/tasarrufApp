package repository

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/offer"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

// OfferRepository implements offer Repository interface
type OfferRepository struct {
	DB *gorm.DB
}

// CreateOfferRepository returns a new instance of offer Repository interface
func CreateOfferRepository(db *gorm.DB) offer.Repository {
	repo := OfferRepository{
		DB: db,
	}
	return &repo
}

// Create a new offer transaction
func (r *OfferRepository) Create(ctx context.Context, o *entities.Offer) (*entities.Offer, error) {
	dbt := r.DB.Create(o)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error creating offer")
	}
	return o, nil
}

// GetByID returns the offer transaction with the given ID
func (r *OfferRepository) GetByID(ctx context.Context, ID uint) (*entities.Offer, error) {
	var o entities.Offer
	dbt := r.DB.Where("id = ?", ID).Find(&o)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.Wrap(dbt.Error, "error getting offer with the given ID")
	}
	return &o, nil
}

// GetByUser returns the offers for the user with the given ID
func (r *OfferRepository) GetByUser(ctx context.Context, userID uint) ([]entities.Offer, error) {
	var offers []entities.Offer
	dbt := r.DB.Where("customer_id = ?", userID).Find(&offers)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.Wrap(dbt.Error, "error getting offers for the given user")
	}
	return offers, nil
}

// GetByPartner returns the offers for the partner with the given ID
func (r *OfferRepository) GetByPartner(ctx context.Context, partnerID uint, startDate time.Time, endDate time.Time) ([]entities.Offer, error) {
	var offers []entities.Offer
	dbt := r.DB.Where("partner_id = ? AND created_at > ? AND created_at < ?", partnerID, startDate, endDate).Find(&offers)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.Wrap(dbt.Error, "error getting offers for the given partner")
	}
	return offers, nil
}

// GetCountByPartnerAndCustomer returns the count of offers consumed by this partner to this customer for the given time range
func (r *OfferRepository) GetCountByPartnerAndCustomer(ctx context.Context, partnerID uint, customerID uint, startDate time.Time, endDate time.Time) (int, error) {
	var offers []entities.Offer
	dbt := r.DB.Where("partner_id = ? AND customer_id = ? AND created_at > ? AND created_at < ?", partnerID, customerID, startDate, endDate).Find(&offers)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return 0, nil
		}
		return 0, errors.Wrap(dbt.Error, "error getting offers for the given partner")
	}
	return len(offers), nil
}

// GetOffersCount returns the count of offers
func (r *OfferRepository) GetOffersCount(ctx context.Context) (int, error) {
	var offers []entities.Offer
	var count int
	dbt := r.DB.Find(&offers).Count(&count)
	if dbt.Error != nil {
		return 0, errors.Wrap(dbt.Error, "error getting offers count")
	}
	return count, nil
}

// GetAllOffers returns the count of offers
func (r *OfferRepository) GetAllOffers(ctx context.Context) ([]entities.Offer, error) {
	var offers []entities.Offer
	dbt := r.DB.Find(&offers)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting offers count")
	}
	return offers, nil
}
