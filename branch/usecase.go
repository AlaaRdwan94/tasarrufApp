package branch

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
)

// Usecase implements branch business logic
type Usecase interface {
	Create(ctx context.Context, b *entities.Branch) (*entities.Branch, error)
	Delete(ctx context.Context, b *entities.Branch) (*entities.Branch, error)
	Edit(ctx context.Context, b *entities.Branch) (*entities.Branch, error)
	GetByID(ctx context.Context, branchID uint) (*entities.Branch, error)
	GetByOwner(ctx context.Context, ownerID uint) ([]entities.Branch, error)
	GetByCategory(ctx context.Context, categoryID uint) ([]entities.Branch, error)
	GetByLocation(ctx context.Context, country string, cityID uint) ([]entities.Branch, error)
	SearchBranches(ctx context.Context, cityID uint, categoryID uint, brandName string) ([]entities.Branch, error)
	CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error)
	DeleteCategory(ctx context.Context, categoryID uint) (*entities.Category, error)
	GetCategories(ctx context.Context) ([]entities.Category, error)
	CreateCity(ctx context.Context, city *entities.City) (*entities.City, error)
	UpdateCity(ctx context.Context, city *entities.City) (*entities.City, error)
	DeleteCity(ctx context.Context, cityID uint) (*entities.City, error)
	GetCityByID(ctx context.Context, cityID uint) (*entities.City, error)
	GetAllCities(ctx context.Context) ([]entities.City, error)
	EditCategory(ctx context.Context, category *entities.Category) (*entities.Category, error)
}
