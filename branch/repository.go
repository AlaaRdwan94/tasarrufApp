package branch

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
)

// Repository defines the repository contract of methods related to the branch data model
type Repository interface {
	Create(ctx context.Context, branch *entities.Branch) (*entities.Branch, error)
	Delete(ctx context.Context, branch *entities.Branch) (*entities.Branch, error)
	Edit(ctx context.Context, branch *entities.Branch) (*entities.Branch, error)
	GetByID(ctx context.Context, branchID uint) (*entities.Branch, error)
	GetByOwner(ctx context.Context, ownerID uint) ([]entities.Branch, error)
	GetByLocation(ctx context.Context, country string, cityID uint) ([]entities.Branch, error)
	CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error)
	DeleteCategory(ctx context.Context, categoryID uint) (*entities.Category, error)
	GetBranchesByCategory(ctx context.Context, categoryID uint) ([]entities.Branch, error)
	GetCategoryByID(ctx context.Context, categoryID uint) (*entities.Category, error)
	GetCategories(ctx context.Context) ([]entities.Category, error)
	EditCategory(ctx context.Context, category *entities.Category) (*entities.Category, error)
	CreateCity(ctx context.Context, city *entities.City) (*entities.City, error)
	UpdateCity(ctx context.Context, city *entities.City) (*entities.City, error)
	DeleteCity(ctx context.Context, city *entities.City) (*entities.City, error)
	GetCityByID(ctx context.Context, cityID uint) (*entities.City, error)
	GetCities(ctx context.Context) ([]entities.City, error)
	SearchBranches(ctx context.Context, city uint, category uint, name string) ([]entities.Branch, error)
}
