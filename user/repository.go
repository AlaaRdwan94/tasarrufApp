// Copyright 2019 NOVA Solutions Co. All Rights Reserved.
//

package user

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"time"
)

// Repository defines the contract of the user repository
type Repository interface {
	CreateCustomer(ctx context.Context, u *entities.User) (*entities.User, error)
	CreatePartner(ctx context.Context, u *entities.User, profile *entities.PartnerProfile) (*entities.User, error)
	UpdateCustomer(ctx context.Context, u *entities.User) (*entities.User, error)
	UpdatePartner(ctx context.Context, u *entities.User, profile *entities.PartnerProfile) (*entities.User, error)
	GetByID(ctx context.Context, ID uint) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByPhone(ctx context.Context, phone string) (*entities.User, error)
	SoftDelete(ctx context.Context, u *entities.User) (*entities.User, error)
	Delete(ctx context.Context, u *entities.User) (*entities.User, error)
	UpdatePartnerProfile(ctx context.Context, p *entities.PartnerProfile) (*entities.PartnerProfile, error)
	DeletePartnerProfile(ctx context.Context, p *entities.PartnerProfile) (*entities.PartnerProfile, error)
	GetPartnerProfileByUserID(ctx context.Context, userID uint) (*entities.PartnerProfile, error)
	CreatePartnerPhoto(ctx context.Context, photo *entities.PartnerPhoto) (*entities.PartnerPhoto, error)
	DeletePartnerPhoto(ctx context.Context, photoID uint) (*entities.PartnerPhoto, error)
	GetPartnerPhotosByPartnerProfileID(ctx context.Context, partnerID uint) ([]entities.PartnerPhoto, error)
	SaveOTP(ctx context.Context, password *entities.OTP) (*entities.OTP, error)
	GetOTPByUser(ctx context.Context, userID uint) (*entities.OTP, error)
	GetCustomersCount(ctx context.Context) (int, error)
	GetPartnersCount(ctx context.Context) (int, error)
	GetAllParnters(ctx context.Context) ([]entities.User, error)
	GetAllCustomers(ctx context.Context) ([]entities.User, error)
	GetNotApprovedPartners(ctx context.Context) ([]entities.User, error)
	CreateExclusiveRecord(ctx context.Context, partnerID uint) error
	GetExclusiveOffers(ctx context.Context) ([]entities.Exclusive, error)
	GetExclusiveOffer(ctx context.Context, partnerID uint) (*entities.Exclusive, error)
	DeleteExclusiveOffer(ctx context.Context, exclusive *entities.Exclusive) (*entities.Exclusive, error)
	CreateShare(ctx context.Context, share *entities.Share) error
	GetSharesByCustomer(ctx context.Context, customerID uint, startDate time.Time, endDate time.Time) (int, error)
	SearchUsers(ctx context.Context, searchTerm string) ([]entities.User, error)
	GetSharablePartners(ctx context.Context) ([]entities.User, error)
}
