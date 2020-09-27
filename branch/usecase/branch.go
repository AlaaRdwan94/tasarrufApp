package usecase

import (
	"context"

	"github.com/ahmedaabouzied/tasarruf/branch"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/subscription"
	"github.com/ahmedaabouzied/tasarruf/user"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// BranchUsecase implements branch usecase interface
type BranchUsecase struct {
	BranchRepo       branch.Repository
	UserRepo         user.Repository
	SubscriptionRepo subscription.Repository
}

// CreateBranchUsecase returns branch usecase implementation
func CreateBranchUsecase(branchRepo branch.Repository, userRepo user.Repository, subcriptionRepo subscription.Repository) branch.Usecase {
	u := BranchUsecase{
		BranchRepo:       branchRepo,
		UserRepo:         userRepo,
		SubscriptionRepo: subcriptionRepo,
	}
	return &u
}

// Create creates a new branch
func (u *BranchUsecase) Create(ctx context.Context, b *entities.Branch) (*entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	creatorUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting partner")
	}
	if creatorUser.AccountType != "partner" {
		err = errors.New("only parner users are allowed to create branches")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	partnerProfile, err := u.UserRepo.GetPartnerProfileByUserID(ctx, currentUserID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "error getting partner profile")
	}
	b.OwnerID = creatorUser.ID
	b.Owner = creatorUser
	b.CategoryID = partnerProfile.CategoryID
	createdBranch, err := u.BranchRepo.Create(ctx, b)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "respository error while creating branch")
	}
	city, err := u.BranchRepo.GetCityByID(ctx, createdBranch.CityID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "respository error while getting city")
	}
	createdBranch.City = *city
	cancelFunc()
	return createdBranch, nil
}

// Delete soft deletes the given branch
func (u *BranchUsecase) Delete(ctx context.Context, b *entities.Branch) (*entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	if b.OwnerID != currentUserID {
		err := errors.New("user is not the owner of this branch")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	deletedBranch, err := u.BranchRepo.Delete(ctx, b)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "respository error while deleting branch")
	}
	city, err := u.BranchRepo.GetCityByID(ctx, deletedBranch.CityID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "respository error while getting city")
	}
	deletedBranch.City = *city
	cancelFunc()
	return deletedBranch, nil
}

// Edit updates the new branch database record
func (u *BranchUsecase) Edit(ctx context.Context, b *entities.Branch) (*entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	toBeEditedBranch, err := u.BranchRepo.GetByID(ctx, b.ID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while feteching branch")
	}
	if currentUserID != toBeEditedBranch.OwnerID {
		err = errors.New("current user is not the owner of this branch")
		log.Error(currentUserID, toBeEditedBranch.ID, err)
		cancelFunc()
		return nil, err
	}
	partnerProfile, err := u.UserRepo.GetPartnerProfileByUserID(ctx, currentUserID)
	if err != nil {
		log.Error(currentUserID, toBeEditedBranch.ID, err)
		cancelFunc()
		return nil, err
	}
	toBeEditedBranch.CategoryID = partnerProfile.CategoryID
	city, err := u.BranchRepo.GetCityByID(ctx, b.CityID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "respository error while getting city")
	}
	b.City = *city
	editedBranch, err := u.BranchRepo.Edit(ctx, b)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "respository error while editing branch")
	}
	cancelFunc()
	return editedBranch, nil
}

// GetByID returns the branch with the given ID
func (u *BranchUsecase) GetByID(ctx context.Context, branchID uint) (*entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	fetchedBranch, err := u.BranchRepo.GetByID(ctx, branchID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "respository error while fetching branch")
	}
	city, err := u.BranchRepo.GetCityByID(ctx, fetchedBranch.CityID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "respository error while getting city")
	}
	fetchedBranch.City = *city
	cancelFunc()
	return fetchedBranch, nil
}

// GetByOwner returns the branches associated with the given owner.
func (u *BranchUsecase) GetByOwner(ctx context.Context, ownerID uint) ([]entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	fetchedBranches, err := u.BranchRepo.GetByOwner(ctx, ownerID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "respository error while fetching branches of the given owner")
	}
	cancelFunc()
	return fetchedBranches, nil
}

// GetByLocation returns the branches associated with the given country , city
func (u *BranchUsecase) GetByLocation(ctx context.Context, country string, cityID uint) ([]entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	fetchedBranches, err := u.BranchRepo.GetByLocation(ctx, country, cityID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "respository error while fetching branches of the given owner")
	}
	for i, branch := range fetchedBranches {
		partner, err := u.UserRepo.GetByID(ctx, branch.OwnerID)
		if err != nil {
			log.Error(err)
			cancelFunc()
			return nil, errors.Wrap(err, "respository error while fetching branches of the given owner")
		}
		fetchedBranches[i].Owner = partner
	}
	cancelFunc()
	return fetchedBranches, nil
}

// CreateCategory creates a new category
func (u *BranchUsecase) CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while getting user")
	}
	if !currentUser.IsAdmin() {
		err = errors.New("only admins are allowed to create categories")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	category, err = u.BranchRepo.CreateCategory(ctx, category)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while getting user")
	}
	cancelFunc()
	return category, nil
}

// EditCategory modifies the category fields.
func (u *BranchUsecase) EditCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while getting user")
	}
	if !currentUser.IsAdmin() {
		err = errors.New("only admins are allowed to edit categories")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	category, err = u.BranchRepo.EditCategory(ctx, category)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while getting user")
	}
	cancelFunc()
	return category, nil

}

func (u *BranchUsecase) DeleteCategory(ctx context.Context, categoryID uint) (*entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while getting user")
	}
	if !currentUser.IsAdmin() {
		err = errors.New("only admins are allowed to delete categories")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	category, err := u.BranchRepo.DeleteCategory(ctx, categoryID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while getting user")
	}
	cancelFunc()
	return category, nil
}

func (u *BranchUsecase) GetByCategory(ctx context.Context, categoryID uint) ([]entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	branches, err := u.BranchRepo.GetBranchesByCategory(ctx, categoryID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while getting categories")
	}
	for i, branch := range branches {
		owner, err := u.UserRepo.GetByID(ctx, branch.OwnerID)
		if err != nil {
			log.Error(err)
			cancelFunc()
			return nil, errors.Wrap(err, "repository error while getting owner")
		}
		branches[i].Owner = owner
	}
	cancelFunc()
	return branches, nil
}

// GetCategories returns all categories
func (u *BranchUsecase) GetCategories(ctx context.Context) ([]entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while getting user subscription")
	}
	isSubscribed := true
	subscription, err := u.SubscriptionRepo.GetSubscriptionByUser(ctx, currentUserID)
	if err != nil {
		isSubscribed = false
	}
	if subscription == nil {
		isSubscribed = false
	}
	if currentUser.IsAdmin() || currentUser.IsPartner() || !isSubscribed {
		categories, err := u.BranchRepo.GetCategories(ctx)
		if err != nil {
			log.Error(err)
			cancelFunc()
			return nil, errors.Wrap(err, "repository error while getting categories")
		}
		cancelFunc()
		return categories, nil

	} else {
		subscription, err := u.SubscriptionRepo.GetSubscriptionByUser(ctx, currentUserID)
		if err != nil {
			log.Error(err)
			cancelFunc()
			return nil, errors.Wrap(err, "repository error while getting user subscription")
		}
		categories, err := u.SubscriptionRepo.GetCategoriesByPlanID(ctx, subscription.PlanID)
		if err != nil {
			log.Error(err)
			cancelFunc()
			return nil, errors.Wrap(err, "repository error while getting categories")
		}
		cancelFunc()
		return categories, nil
	}
}

// CreateCity creates a new city
func (u *BranchUsecase) CreateCity(ctx context.Context, city *entities.City) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		err = errors.Wrap(err, "error getting current user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if !currentUser.IsAdmin() {
		err = errors.New("not authorized")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	city, err = u.BranchRepo.CreateCity(ctx, city)
	if err != nil {
		err = errors.Wrap(err, "error creating city")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return city, err
}

// DeleteCity deletes a city
func (u *BranchUsecase) DeleteCity(ctx context.Context, cityID uint) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		err = errors.Wrap(err, "error getting current user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if !currentUser.IsAdmin() {
		err = errors.New("not authorized")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	city, err := u.BranchRepo.GetCityByID(ctx, cityID)
	if err != nil {
		err = errors.Wrap(err, "error getting city")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	city, err = u.BranchRepo.DeleteCity(ctx, city)
	if err != nil {
		err = errors.Wrap(err, "error deleting city")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return city, err
}

// UpdateCity updates the city
func (u *BranchUsecase) UpdateCity(ctx context.Context, updatedCity *entities.City) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		err = errors.Wrap(err, "error getting current user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	city, err := u.BranchRepo.GetCityByID(ctx, updatedCity.ID)
	if err != nil {
		err = errors.Wrap(err, "error getting city")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	err = city.SetEnglishName(currentUser, updatedCity.EnglishName)
	if err != nil {
		err = errors.Wrap(err, "error updating city englishName")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	err = city.SetTurkishName(currentUser, updatedCity.TurkishName)
	if err != nil {
		err = errors.Wrap(err, "error updating city turkishName")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	city, err = u.BranchRepo.UpdateCity(ctx, city)
	if err != nil {
		err = errors.Wrap(err, "error updating city")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return city, nil
}

// GetCityByID returns the city with the given ID
func (u *BranchUsecase) GetCityByID(ctx context.Context, cityID uint) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	city, err := u.BranchRepo.GetCityByID(ctx, cityID)
	if err != nil {
		err = errors.Wrap(err, "error getting city")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return city, nil
}

// GetAllCities returns all the cities
func (u *BranchUsecase) GetAllCities(ctx context.Context) ([]entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	cities, err := u.BranchRepo.GetCities(ctx)
	if err != nil {
		err = errors.Wrap(err, "error getting cities")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return cities, nil
}

// SearchBranches searches the branches with brandName , cityID, categoryID
func (u *BranchUsecase) SearchBranches(ctx context.Context, cityID uint, categoryID uint, brandName string) ([]entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	branches, err := u.BranchRepo.SearchBranches(ctx, cityID, categoryID, brandName)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	for i, branch := range branches {
		log.Info("owner id = ", branch.OwnerID)
		owner, err := u.UserRepo.GetByID(ctx, branch.OwnerID)
		if err != nil {
			cancelFunc()
			return nil, errors.Wrap(err, "error getting branch owner")
		}
		branches[i].Owner = owner
		city, err := u.BranchRepo.GetCityByID(ctx, branch.CityID)
		if err != nil {
			cancelFunc()
			return nil, errors.Wrap(err, "error getting city")
		}
		branches[i].City = *city
		ownerCity, err := u.BranchRepo.GetCityByID(ctx, branches[i].Owner.CityID)
		if err != nil {
			cancelFunc()
			return nil, errors.Wrap(err, "error getting city")
		}
		branches[i].Owner.City = *ownerCity
		branches[i].Owner.PartnerProfile.City = *ownerCity
	}
	cancelFunc()
	return branches, nil
}
