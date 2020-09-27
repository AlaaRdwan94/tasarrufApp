package support

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
)

// Repository represents the Support Info repository
type Repository interface {
	Create(ctx context.Context, info *entities.SupportInfo) (*entities.SupportInfo, error)
	Update(ctx context.Context, info *entities.SupportInfo) (*entities.SupportInfo, error)
	Delete(ctx context.Context, info *entities.SupportInfo) (*entities.SupportInfo, error)
	GetSupportInfo(ctx context.Context) (*entities.SupportInfo, error)
}
