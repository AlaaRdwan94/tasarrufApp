package support

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
)

// Usecase represents the Support Info usecase contract
type Usecase interface {
	Create(ctx context.Context, info *entities.SupportInfo) (*entities.SupportInfo, error)
	Update(ctx context.Context, info *entities.SupportInfo) (*entities.SupportInfo, error)
	GetSupportInfo(ctx context.Context) (*entities.SupportInfo, error)
}
