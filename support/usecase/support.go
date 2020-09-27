package usecase

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/support"
	"github.com/ahmedaabouzied/tasarruf/user"
	"github.com/pkg/errors"
)

// SupportUsecase represents an implementaion of support usecase interface
type SupportUsecase struct {
	SupportRepo support.Repository
	UserRepo    user.Repository
}

func CreateSupportUsecase(supportRepo support.Repository, userRepo user.Repository) support.Usecase {
	u := SupportUsecase{
		SupportRepo: supportRepo,
		UserRepo:    userRepo,
	}
	return &u
}

// Create creates a new support infor record
func (u *SupportUsecase) Create(ctx context.Context, info *entities.SupportInfo) (*entities.SupportInfo, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		cancelFunc()
		return nil, errors.New("user is not authorized to create support record")
	}
	info, err = u.SupportRepo.Create(ctx, info)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error creating support record")
	}
	cancelFunc()
	return info, nil
}

// Update updates the support info record
func (u *SupportUsecase) Update(ctx context.Context, info *entities.SupportInfo) (*entities.SupportInfo, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.UserRepo.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		cancelFunc()
		return nil, errors.New("user is not authorized to update support record")
	}
	toUpdateInfo, err := u.SupportRepo.GetSupportInfo(ctx)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error updating support info")
	}
	toUpdateInfo.Email = info.Email
	toUpdateInfo.Mobile = info.Mobile
	info, err = u.SupportRepo.Update(ctx, toUpdateInfo)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error updating support record")
	}
	cancelFunc()
	return info, nil
}

// GetSupportInfo returns the support info record
func (u *SupportUsecase) GetSupportInfo(ctx context.Context) (*entities.SupportInfo, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	info, err := u.SupportRepo.GetSupportInfo(ctx)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error updating support record")
	}
	cancelFunc()
	return info, nil
}
