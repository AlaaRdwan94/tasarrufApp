package user

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"mime/multipart"
)

// Usecase represents the user business logic
type Usecase interface {
	CreateCustomer(ctx context.Context, u *entities.User, newPassword string) (*entities.User, error)
	CreatePartner(ctx context.Context, u *entities.User, profile *entities.PartnerProfile, newPassword string) (*entities.User, error)
	IsEmailRegistered(ctx context.Context, email string) (bool, error)
	IsPhoneRegistered(ctx context.Context, phone string) (bool, error)
	EmailLogin(ctx context.Context, email string, password string) (*entities.User, string, error)
	PhoneLogin(ctx context.Context, phone string, password string) (*entities.User, string, error)
	GetUser(ctx context.Context, ID uint) (*entities.User, error)
	DeleteUser(ctx context.Context, ID uint) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	RecoverPassword(ctx context.Context, mobile string) error
	AdminRecoverPassword(ctx context.Context, mobile string) error
	UpdatePassword(ctx context.Context, newPassword string) (*entities.User, error)
	UpdateProfileImage(ctx context.Context, fileHeader *multipart.FileHeader) (*entities.User, error)
	UpdateTradeLicense(ctx context.Context, fileHeader *multipart.FileHeader) (*entities.Partner, error)
	AddPartnerPhoto(ctx context.Context, fileHeader *multipart.FileHeader) (*entities.User, error)
	DeletePartnerPhoto(ctx context.Context, photoID uint) (*entities.User, error)
	UpdatePartnerProfile(ctx context.Context, profile *entities.PartnerProfile) (*entities.User, error)
	VerifyUser(ctx context.Context, code string) (*entities.User, error)
	ResendVerficationCode(ctx context.Context) error
	ValidateCustomerPartnerIntegrity(ctx context.Context, customerID uint, partnerID uint) (*entities.User, *entities.Subscription, error)
	GetCustomersCount(ctx context.Context) (int, error)
	GetCustomerByID(ctx context.Context, ID uint) (*entities.Customer, error)
	GetPartnersCount(ctx context.Context) (int, error)
	GetPartnerByID(ctx context.Context, ID uint) (*entities.Partner, error)
	GetAllCustomers(ctx context.Context) ([]entities.User, error)
	GetAllParnters(ctx context.Context) ([]entities.Partner, error)
	GetNotApprovedPartners(ctx context.Context) ([]entities.User, error)
	ApprovePartner(ctx context.Context, partnerID uint) (*entities.User, error)
	SetPartnerAsExclusive(ctx context.Context, partnerID uint) (*entities.Partner, error)
	RemovePartnerAsExclusive(ctx context.Context, partnerID uint) (*entities.Partner, error)
	GetExclusivePartners(ctx context.Context) ([]entities.Partner, error)
	Share(ctx context.Context) error
	AdminDeleteUser(ctx context.Context, userID uint) (*entities.User, error)
	ToggleIsSharable(ctx context.Context, partnerID uint) error
	SearchUsers(ctx context.Context, searchTerm string) ([]entities.IUser, error)
	ToggleActiveProperty(ctx context.Context, userID uint) (entities.IUser, error)
}
