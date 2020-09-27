package repository

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/review"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ReviewRepository implements the review Repository interface
type ReviewRepository struct {
	DB *gorm.DB
}

// CreateReviewRepository returns an instance of the review repository interface
func CreateReviewRepository(db *gorm.DB) review.Repository {
	r := &ReviewRepository{
		DB: db,
	}
	return r
}

// Create creates a new review DB record
func (r *ReviewRepository) Create(ctx context.Context, review *entities.Review) (*entities.Review, error) {
	dbt := r.DB.Create(review)
	if dbt.Error != nil {
		ctx.Done()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error creating review")
	}
	return review, nil
}

// Update updates the db record of the given review
func (r *ReviewRepository) Update(ctx context.Context, review *entities.Review) (*entities.Review, error) {
	dbt := r.DB.Save(review)
	if dbt.Error != nil {
		ctx.Done()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error updating reveiew")
	}
	return review, nil
}

// Delete soft deletes the db record of the given review
func (r *ReviewRepository) Delete(ctx context.Context, review *entities.Review) (*entities.Review, error) {
	dbt := r.DB.Delete(review)
	if dbt.Error != nil {
		ctx.Done()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error deleting review")
	}
	return review, nil

}

// GetByID returns the db record of the review with the given ID
func (r *ReviewRepository) GetByID(ctx context.Context, ID uint) (*entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var review entities.Review
	dbt := r.DB.Where("id = ?", ID).Find(&review)
	if dbt.RecordNotFound() {
		cancelFunc()
		return nil, nil
	}
	if dbt.Error != nil {
		cancelFunc()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error getting review")
	}
	cancelFunc()
	return &review, nil
}

// GetByCustomerID returns an array of records of reviews with the given customer ID.
func (r *ReviewRepository) GetByCustomerID(ctx context.Context, ID uint) ([]entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var reviews []entities.Review
	dbt := r.DB.Where("customer_id = ?", ID).Find(&reviews)
	if dbt.RecordNotFound() {
		cancelFunc()
		return nil, nil
	}
	if dbt.Error != nil {
		cancelFunc()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error getting reviews")
	}
	cancelFunc()
	return reviews, nil
}

// GetByPartnerID returns an array of records of reviews with the given customer ID.
func (r *ReviewRepository) GetByPartnerID(ctx context.Context, ID uint) ([]entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var reviews []entities.Review
	dbt := r.DB.Where("partner_id = ?", ID).Find(&reviews)
	if dbt.RecordNotFound() {
		cancelFunc()
		return nil, nil
	}
	if dbt.Error != nil {
		cancelFunc()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error getting reviews")
	}
	cancelFunc()
	return reviews, nil
}

func (r *ReviewRepository) GetAverageRatings(ctx context.Context, partnerID uint) (float64, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var avr float64
	row := r.DB.Table("reviews").Select("avg(stars)").Where("partner_id = ?", partnerID).Row()
	err := row.Scan(&avr)
	if err != nil {
		cancelFunc()
		return 0, err
	}
	cancelFunc()
	return avr, nil
}
