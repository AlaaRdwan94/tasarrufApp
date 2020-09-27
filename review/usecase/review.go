package usecase

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/branch"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/review"
	"github.com/ahmedaabouzied/tasarruf/user"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ReviewUsecase implements review usecase interface
type ReviewUsecase struct {
	BranchRepo branch.Repository
	UserRepo   user.Repository
	ReviewRepo review.Repository
}

// CreateReviewUsecase returns review usecase implementation
func CreateReviewUsecase(reviewRepo review.Repository, userRepo user.Repository, branchRepo branch.Repository) review.Usecase {
	u := ReviewUsecase{
		BranchRepo: branchRepo,
		UserRepo:   userRepo,
		ReviewRepo: reviewRepo,
	}
	return &u
}

func (u *ReviewUsecase) Create(ctx context.Context, review *entities.Review) (*entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err

	}
	if currentUser.AccountType != "user" {
		err := errors.New("only customer users are allowed to create reviews")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	partner, err := u.UserRepo.GetByID(ctx, review.PartnerID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err

	}
	if partner.AccountType != "partner" {
		err := errors.New("reviews are only allowed on partner profiles")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	review.CustomerID = currentUser.ID
	review.PartnerID = partner.ID
	review, err = u.ReviewRepo.Create(ctx, review)
	if err != nil {
		err = errors.Wrap(err, "repsitory error while creating review")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	review.Partner = *partner
	review.Customer = *currentUser
	cancelFunc()
	return review, nil
}

// Update updates the given review
func (u *ReviewUsecase) Update(ctx context.Context, review *entities.Review) (*entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	if review.CustomerID != currentUserID {
		err := errors.New("current user is not the owner of the review")
		log.Error(err)
		cancelFunc()
		return nil, err

	}
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err

	}
	partner, err := u.UserRepo.GetByID(ctx, review.PartnerID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	review, err = u.ReviewRepo.Update(ctx, review)
	if err != nil {
		err = errors.Wrap(err, "repsitory error while updating review")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	review.Partner = *partner
	review.Customer = *currentUser
	cancelFunc()
	return review, nil
}

// Delete deletes the given review
func (u *ReviewUsecase) Delete(ctx context.Context, review *entities.Review) (*entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	if review.CustomerID != currentUserID {
		err := errors.New("current user is not the owner of the review")
		log.Error(err)
		cancelFunc()
		return nil, err

	}
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err

	}
	partner, err := u.UserRepo.GetByID(ctx, review.PartnerID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	review, err = u.ReviewRepo.Delete(ctx, review)
	if err != nil {
		err = errors.Wrap(err, "repsitory error while deleting review")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	review.Partner = *partner
	review.Customer = *currentUser
	cancelFunc()
	return review, nil
}

// GetByID returns the review with the given ID
func (u *ReviewUsecase) GetByID(ctx context.Context, ID uint) (*entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	review, err := u.ReviewRepo.GetByID(ctx, ID)
	if err != nil {
		err = errors.Wrap(err, "repsitory error while deleting review")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	customer, err := u.UserRepo.GetByID(ctx, review.CustomerID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err

	}
	partner, err := u.UserRepo.GetByID(ctx, review.PartnerID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	review.Partner = *partner
	review.Customer = *customer
	cancelFunc()
	return review, nil
}

// GetByCustomerID returns an array of reviews by the given customer ID
func (u *ReviewUsecase) GetByCustomerID(ctx context.Context, ID uint) ([]entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	reviews, err := u.ReviewRepo.GetByCustomerID(ctx, ID)
	if err != nil {
		err = errors.Wrap(err, "repsitory error while deleting review")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	for i, review := range reviews {
		customer, err := u.UserRepo.GetByID(ctx, review.CustomerID)
		if err != nil {
			err := errors.Wrap(err, "repository error while getting user")
			log.Error(err)
			cancelFunc()
			return nil, err

		}
		partner, err := u.UserRepo.GetByID(ctx, review.PartnerID)
		if err != nil {
			err := errors.Wrap(err, "repository error while getting user")
			log.Error(err)
			cancelFunc()
			return nil, err
		}
		reviews[i].Partner = *partner
		reviews[i].Customer = *customer
	}
	cancelFunc()
	return reviews, nil
}

// GetByPartnerID returns an array of reviews by the given customer ID
func (u *ReviewUsecase) GetByPartnerID(ctx context.Context, ID uint) ([]entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	reviews, err := u.ReviewRepo.GetByPartnerID(ctx, ID)
	if err != nil {
		err = errors.Wrap(err, "repsitory error while deleting review")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	for i, review := range reviews {
		customer, err := u.UserRepo.GetByID(ctx, review.CustomerID)
		if err != nil {
			err := errors.Wrap(err, "repository error while getting user")
			log.Error(err)
			cancelFunc()
			return nil, err

		}
		partner, err := u.UserRepo.GetByID(ctx, review.PartnerID)
		if err != nil {
			err := errors.Wrap(err, "repository error while getting user")
			log.Error(err)
			cancelFunc()
			return nil, err
		}
		reviews[i].Partner = *partner
		reviews[i].Customer = *customer
	}
	cancelFunc()
	return reviews, nil
}

// GetAverageRatings returns the AVG of stars of the given partner
func (u *ReviewUsecase) GetAverageRatings(ctx context.Context, partnerID uint) (float64, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	average, err := u.ReviewRepo.GetAverageRatings(ctx, partnerID)
	if err != nil {
		err = errors.Wrap(err, "repsitory error while getting average rating")
		log.Error(err)
		cancelFunc()
		return 0, err
	}
	cancelFunc()
	return average, nil
}
