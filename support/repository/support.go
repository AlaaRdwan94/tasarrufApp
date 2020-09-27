package repository

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/support"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// SupportRepository represents a support repository
type SupportRepository struct {
	DB *gorm.DB
}

// CreateSupportRepository creates a new support info repository instance
func CreateSupportRepository(db *gorm.DB) support.Repository {
	repo := SupportRepository{
		DB: db,
	}
	return &repo
}

// Create creates a new Support info record
func (r *SupportRepository) Create(ctx context.Context, info *entities.SupportInfo) (*entities.SupportInfo, error) {
	dbt := r.DB.Create(info)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error creating support info record")
	}
	return info, nil
}

// Update updates the given support info record
func (r *SupportRepository) Update(ctx context.Context, info *entities.SupportInfo) (*entities.SupportInfo, error) {
	dbt := r.DB.Save(info)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error updating support info record")
	}
	return info, nil
}

// Delete deletes the given support info record
func (r *SupportRepository) Delete(ctx context.Context, info *entities.SupportInfo) (*entities.SupportInfo, error) {
	dbt := r.DB.Delete(info)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error deleting support info record")
	}
	return info, nil
}

// GetSupportInfo returns the top support info record
func (r *SupportRepository) GetSupportInfo(ctx context.Context) (*entities.SupportInfo, error) {
	var info entities.SupportInfo
	dbt := r.DB.Order("created_at DESC").First(&info)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting support information")
	}
	return &info, nil
}
