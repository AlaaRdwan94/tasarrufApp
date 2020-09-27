package repository

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/branch"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// BranchRepository implements branch repository interface
type BranchRepository struct {
	DB *gorm.DB
}

// CreateBranchRepository returns a branch repository struct
func CreateBranchRepository(db *gorm.DB) branch.Repository {
	repo := BranchRepository{
		DB: db,
	}
	return &repo
}

// Create creates a new branch
func (r *BranchRepository) Create(ctx context.Context, branch *entities.Branch) (*entities.Branch, error) {
	dbt := r.DB.Create(branch)
	if dbt.Error != nil {
		ctx.Done()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error creating branch")
	}
	return branch, nil
}

// Delete deletes the given branch
func (r *BranchRepository) Delete(ctx context.Context, branch *entities.Branch) (*entities.Branch, error) {
	dbt := r.DB.Delete(branch)
	if dbt.Error != nil {
		ctx.Done()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error deleting branch")
	}
	return branch, nil
}

// Edit saves the given branch
func (r *BranchRepository) Edit(ctx context.Context, branch *entities.Branch) (*entities.Branch, error) {
	log.Info(branch.CityID)
	r.DB.LogMode(true)
	r.DB.Debug()
	dbt := r.DB.Save(branch)
	if dbt.Error != nil {
		ctx.Done()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error editing branch")
	}
	return branch, nil

}

// GetByID returns the branch with the given ID
func (r *BranchRepository) GetByID(ctx context.Context, branchID uint) (*entities.Branch, error) {
	var b entities.Branch
	dbt := r.DB.Where("id = ?", branchID).Find(&b)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			ctx.Done()
			return nil, nil
		}
		ctx.Done()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error getting branch with the given id")
	}
	city, err := r.GetCityByID(ctx, b.CityID)
	if err != nil {
		ctx.Done()
		log.Error(err)
		return nil, err
	}
	b.City = *city
	return &b, nil
}

// GetByOwner returns the list of branches associated with the given owner ID
func (r *BranchRepository) GetByOwner(ctx context.Context, ownerID uint) ([]entities.Branch, error) {
	var b []entities.Branch
	dbt := r.DB.Where("owner_id = ?", ownerID).Find(&b)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			ctx.Done()
			return nil, nil
		}
		ctx.Done()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error getting branches associated with the given owner")
	}
	for id, branch := range b {
		city, err := r.GetCityByID(ctx, branch.CityID)
		if err != nil {
			ctx.Done()
			log.Error(err)
			return nil, err
		}
		b[id].City = *city
	}
	return b, nil
}

// GetByLocation returns the list of branches with the given city and country
func (r *BranchRepository) GetByLocation(ctx context.Context, country string, cityID uint) ([]entities.Branch, error) {
	var b []entities.Branch
	dbt := r.DB.Where("country = ? AND city_id = ?", country, cityID).Find(&b)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			ctx.Done()
			return nil, nil
		}
		ctx.Done()
		log.Error(dbt.Error)
		return nil, errors.Wrap(dbt.Error, "error getting branches in the given location")
	}
	for id, branch := range b {
		city, err := r.GetCityByID(ctx, branch.CityID)
		if err != nil {
			ctx.Done()
			log.Error(err)
			return nil, err
		}
		b[id].City = *city
	}
	return b, nil
}

// CreateCategory creates a DB record of the given category
func (r *BranchRepository) CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	dbt := r.DB.Create(category)
	if dbt.Error != nil {
		log.Error(dbt.Error)
		cancelFunc()
		return nil, dbt.Error
	}
	cancelFunc()
	return category, nil
}

// EditCategory updates the given category to the databse
func (r *BranchRepository) EditCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	dbt := r.DB.Save(category)
	if dbt.Error != nil {
		log.Error(dbt.Error)
		cancelFunc()
		return nil, dbt.Error
	}
	cancelFunc()
	return category, nil
}

// GetCategories returns all the categories.
// If there are none it will return nil.
func (r *BranchRepository) GetCategories(ctx context.Context) ([]entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var categories []entities.Category
	dbt := r.DB.Order("english_name ASC").Find(&categories)
	if dbt.Error != nil {
		log.Error(dbt.Error)
		cancelFunc()
		return nil, dbt.Error
	}
	cancelFunc()
	return categories, nil
}

// DeleteCategory deletes the category with the given ID.
func (r *BranchRepository) DeleteCategory(ctx context.Context, categoryID uint) (*entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	category, err := r.GetCategoryByID(ctx, categoryID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if category == nil {
		err := errors.New("category not found")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	dbt := r.DB.Delete(category)
	if dbt.Error != nil {
		log.Error(dbt.Error)
		cancelFunc()
		return nil, dbt.Error
	}
	cancelFunc()
	return category, nil
}

// GetCategoryByID returns the category with the given ID
func (r *BranchRepository) GetCategoryByID(ctx context.Context, categoryID uint) (*entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var category entities.Category
	dbt := r.DB.Where("id = ?", categoryID).Find(&category)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			cancelFunc()
			return nil, nil
		}
		log.Error(dbt.Error)
		cancelFunc()
		return nil, dbt.Error
	}
	cancelFunc()
	return &category, nil
}

// GetBranchesByCategory returns an array of branches associated with the category with the given ID.
func (r *BranchRepository) GetBranchesByCategory(ctx context.Context, categoryID uint) ([]entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var branches []entities.Branch
	dbt := r.DB.Where("category_id = ?", categoryID).Find(&branches)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			cancelFunc()
			return nil, nil
		}
		log.Error(dbt.Error)
		cancelFunc()
		return nil, dbt.Error
	}
	cancelFunc()
	return branches, nil
}

// CreateCity creates a new city
func (r *BranchRepository) CreateCity(ctx context.Context, city *entities.City) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	dbt := r.DB.Create(city)
	if dbt.Error != nil {
		cancelFunc()
		return nil, dbt.Error
	}
	cancelFunc()
	return city, nil
}

// DeleteCity deletes the given city
func (r *BranchRepository) DeleteCity(ctx context.Context, city *entities.City) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	// get users by this city
	var users []entities.User
	dbt := r.DB.Where("city_id = ?", city.ID).Find(&users)
	// get branches by this city
	var branches []entities.Branch
	dbt = r.DB.Where("city_id = ?", city.ID).Find(&branches)
	if len(users) < 0 || len(branches) > 0 {
		cancelFunc()
		return nil, errors.New("city cannot be deleted")
	}
	dbt = r.DB.Delete(city)
	if dbt.Error != nil {
		cancelFunc()
		return nil, dbt.Error
	}
	cancelFunc()
	return city, nil
}

// UpdateCity saves the given city
func (r *BranchRepository) UpdateCity(ctx context.Context, city *entities.City) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	dbt := r.DB.Save(city)
	if dbt.Error != nil {
		cancelFunc()
		return nil, dbt.Error
	}
	cancelFunc()
	return city, nil

}

// GetCityByID gets the city with the given ID
func (r *BranchRepository) GetCityByID(ctx context.Context, cityID uint) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var city entities.City
	dbt := r.DB.Where("id = ?", cityID).Find(&city)
	if dbt.Error != nil {
		cancelFunc()
		return nil, dbt.Error
	}
	cancelFunc()
	return &city, nil
}

// GetCities returns all the citites
func (r *BranchRepository) GetCities(ctx context.Context) ([]entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var cities []entities.City
	dbt := r.DB.Order("english_name ASC").Find(&cities)
	if dbt.Error != nil {
		cancelFunc()
		return nil, dbt.Error
	}
	cancelFunc()
	return cities, nil
}

// SearchBranches searches branches with brand name , category id , city id
func (r *BranchRepository) SearchBranches(ctx context.Context, city uint, category uint, name string) ([]entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	log.Info("Name = ", name)
	name = "%" + name + "%"
	log.Info(name)
	var branches []entities.Branch
	var dbt *gorm.DB
	r.DB.Debug()
	r.DB.LogMode(true)
	if category > 0 {
		dbt = r.DB.Raw(`
            SELECT * FROM branches WHERE deleted_at IS NULL AND city_id = ? AND category_id = ? AND owner_id IN (
            SELECT partner_id FROM partner_profiles WHERE LOWER(brand_name) LIKE ? AND partner_id IN (SELECT id FROM users WHERE active = true))`, city, category, name).Scan(&branches)
	} else {
		dbt = r.DB.Raw(`
        SELECT * FROM branches WHERE deleted_at IS NULL AND city_id = ? AND owner_id IN (
            SELECT partner_id FROM partner_profiles WHERE LOWER(brand_name) LIKE ? AND partner_id IN (SELECT id FROM users WHERE active = true))`, city, name).Scan(&branches)
	}

	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			cancelFunc()
			return nil, nil
		}
		cancelFunc()
		return nil, errors.Wrap(dbt.Error, "error getting search results")
	}
	cancelFunc()
	return branches, nil
}
