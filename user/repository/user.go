package repository

import (
	"context"

	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/user"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

// UserRepository implements user repository interface
type UserRepository struct {
	DB *gorm.DB
}

// CreateUserRepository returns a user repository struct
func CreateUserRepository(db *gorm.DB) user.Repository {
	repo := UserRepository{
		DB: db,
	}
	return &repo
}

// CreateCustomer creates a new customer user
func (r *UserRepository) CreateCustomer(ctx context.Context, u *entities.User) (*entities.User, error) {
	dbt := r.DB.Create(u)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error creating user")
	}
	return u, nil
}

// CreatePartner creates new partner user , parnter profile db records
func (r *UserRepository) CreatePartner(ctx context.Context, u *entities.User, profile *entities.PartnerProfile) (*entities.User, error) {
	// create partner user
	dbt := r.DB.Create(u)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error creating partner")
	}
	profile.PartnerID = u.ID
	dbt = r.DB.Create(profile)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error creating partner")
	}
	u.PartnerProfile = *profile
	return u, nil
}

// UpdateCustomer updates the customer user details
func (r *UserRepository) UpdateCustomer(ctx context.Context, u *entities.User) (*entities.User, error) {
	dbt := r.DB.Save(u)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error saving user")
	}
	return u, nil
}

// UpdatePartner updates the partner user details
func (r *UserRepository) UpdatePartner(ctx context.Context, u *entities.User, profile *entities.PartnerProfile) (*entities.User, error) {
	dbt := r.DB.Save(u)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error saving partner")
	}
	return u, nil
}

// GetByID returns the user with the given ID
func (r *UserRepository) GetByID(ctx context.Context, ID uint) (*entities.User, error) {
	var u entities.User
	dbt := r.DB.Where("id = ?", ID).Find(&u)
	if dbt.RecordNotFound() {
		return nil, errors.Wrap(dbt.Error, "user not found")
	}
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting user with the given ID")
	}
	if u.AccountType == "partner" {
		partnerProfile, err := r.GetPartnerProfileByUserID(ctx, u.ID)
		if err != nil {
			return nil, errors.Wrap(dbt.Error, "error getting user with the given ID")
		}
		u.PartnerProfile = *partnerProfile
	}

	return &u, nil
}

// UpdatePartnerProfile updates the partner profile of the partner user
func (r *UserRepository) UpdatePartnerProfile(ctx context.Context, p *entities.PartnerProfile) (*entities.PartnerProfile, error) {
	dbt := r.DB.Save(p)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error saving user")
	}
	return p, nil
}

// DeletePartnerProfile Deletes the partner profile of the partner user
func (r *UserRepository) DeletePartnerProfile(ctx context.Context, p *entities.PartnerProfile) (*entities.PartnerProfile, error) {
	dbt := r.DB.Delete(p)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error saving user")
	}
	return p, nil
}

// GetPartnerProfileByUserID returns the partner profile associated with the given partner user ID
func (r *UserRepository) GetPartnerProfileByUserID(ctx context.Context, userID uint) (*entities.PartnerProfile, error) {
	var u entities.PartnerProfile
	dbt := r.DB.Where("partner_id = ?", userID).Find(&u)
	if dbt.RecordNotFound() {
		return nil, errors.Wrap(dbt.Error, "profile not found")
	}
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting profile with the given user ID")
	}
	photos, err := r.GetPartnerPhotosByPartnerProfileID(ctx, u.ID)
	if err != nil {
		return nil, errors.Wrap(dbt.Error, "error getting photos of partner")
	}
	u.PartnerPhotos = photos
	return &u, nil
}

// GetByEmail returns the user with the given email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var u entities.User
	dbt := r.DB.Where("email = ?", email).Find(&u)
	if dbt.RecordNotFound() {
		return nil, errors.Wrap(dbt.Error, "user not found")
	}
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting user with the given ID")
	}
	if u.AccountType == "partner" {
		partnerProfile, err := r.GetPartnerProfileByUserID(ctx, u.ID)
		if err != nil {
			return nil, errors.Wrap(dbt.Error, "error getting user with the given ID")
		}
		u.PartnerProfile = *partnerProfile
	}
	return &u, nil
}

// GetByPhone returns the user with the given email
func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*entities.User, error) {
	var u entities.User
	dbt := r.DB.Where("mobile = ?", phone).Find(&u)
	if dbt.RecordNotFound() {
		return nil, errors.Wrap(dbt.Error, "user not found")
	}
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting user with the given ID")
	}
	return &u, nil
}

// SoftDelete sets the user deleted at
func (r *UserRepository) SoftDelete(ctx context.Context, u *entities.User) (*entities.User, error) {
	dbt := r.DB.Delete(u)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error deleting user")
	}
	return u, nil
}

// Delete the given user
func (r *UserRepository) Delete(ctx context.Context, u *entities.User) (*entities.User, error) {
	dbt := r.DB.Unscoped().Delete(u)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error deleting user")
	}
	return u, nil
}

// CreatePartnerPhoto creates a new DB record of the given partner photo
func (r *UserRepository) CreatePartnerPhoto(ctx context.Context, photo *entities.PartnerPhoto) (*entities.PartnerPhoto, error) {
	dbt := r.DB.Create(photo)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "failed to save partner photo url to DB")
	}
	return photo, nil
}

// DeletePartnerPhoto deletes the given partner photo
func (r *UserRepository) DeletePartnerPhoto(ctx context.Context, photoID uint) (*entities.PartnerPhoto, error) {
	var photo entities.PartnerPhoto
	dbt := r.DB.Where("id = ?", photoID).Find(&photo)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "failed to delete partner photo")
	}
	dbt = r.DB.Delete(&photo)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "failed to delete partner photo")
	}
	return &photo, nil

}

// GetPartnerPhotosByPartnerProfileID returns the array of photos belonging to a partner
func (r *UserRepository) GetPartnerPhotosByPartnerProfileID(ctx context.Context, profileID uint) ([]entities.PartnerPhoto, error) {
	var photos []entities.PartnerPhoto
	dbt := r.DB.Where("partner_profile_ID = ?", profileID).Find(&photos)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.Wrap(dbt.Error, "failed to fetch partner photos")
	}
	return photos, nil
}

// CreateOTP saves the OTP to DB
func (r *UserRepository) SaveOTP(ctx context.Context, password *entities.OTP) (*entities.OTP, error) {
	var oldPassword entities.OTP
	p := &oldPassword
	dbt := r.DB.Where("user_id = ?", password.UserID).Find(p)
	if dbt.RecordNotFound() {
		dbt := r.DB.Create(password)
		if dbt.Error != nil {
			return nil, errors.Wrap(dbt.Error, "error saving OTP to DB")
		}
	}
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting OTP from DB")
	}
	p.HashedPassword = password.HashedPassword
	p.ExpiryDate = password.ExpiryDate
	dbt = r.DB.Save(&oldPassword)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting OTP from DB")
	}
	return &oldPassword, nil
}

// GetOTPByUser returns the OTP associated with the user with the given ID
func (r *UserRepository) GetOTPByUser(ctx context.Context, userID uint) (*entities.OTP, error) {
	var password entities.OTP
	dbt := r.DB.Where("user_id = ?", userID).Find(&password)
	if dbt.RecordNotFound() {
		return nil, nil
	}
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting OTP from DB")
	}
	return &password, nil
}

// GetCustomersCount returns the count of all customer users.
func (r *UserRepository) GetCustomersCount(ctx context.Context) (int, error) {
	var customers []entities.User
	var count int
	dbt := r.DB.Where("account_type = ?", "user").Find(&customers).Count(&count)
	if dbt.Error != nil {
		return 0, errors.Wrap(dbt.Error, "error getting customers count")
	}
	return count, nil
}

// GetPartnersCount returns the count of all partner users.
func (r *UserRepository) GetPartnersCount(ctx context.Context) (int, error) {
	var partners []entities.User
	var count int
	r.DB.LogMode(true)
	r.DB.Debug()
	dbt := r.DB.Where("account_type = ?", "partner").Find(&partners).Count(&count)
	if dbt.Error != nil {
		return 0, errors.Wrap(dbt.Error, "error getting partners count")
	}
	return count, nil
}

// GetAllParnters returns a list of all partner users
func (r *UserRepository) GetAllParnters(ctx context.Context) ([]entities.User, error) {
	var partners []entities.User
	dbt := r.DB.Where("account_type = ?", "partner").Find(&partners)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting partners count")
	}
	return partners, nil
}

// GetAllCustomers returns a list of all partner users
func (r *UserRepository) GetAllCustomers(ctx context.Context) ([]entities.User, error) {
	var customers []entities.User
	dbt := r.DB.Where("account_type = ?", "user").Find(&customers)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting partners count")
	}
	return customers, nil
}

// GetNotApprovedPartners returns a list of partners where approved = false
func (r *UserRepository) GetNotApprovedPartners(ctx context.Context) ([]entities.User, error) {
	var partners []entities.User
	dbt := r.DB.Raw(`SELECT * FROM users WHERE id IN (
        SELECT partner_id FROM partner_profiles WHERE approved = ?
    ) AND deleted_at IS NULL`, false).Scan(&partners)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error getting not approved partners")
	}
	return partners, nil
}

// CreateExclusiveRecord creates a new exclusive partner record
func (r *UserRepository) CreateExclusiveRecord(ctx context.Context, partnerID uint) error {
	exclusive := entities.Exclusive{
		PartnerID: partnerID,
	}
	dbt := r.DB.Create(&exclusive)
	if dbt.Error != nil {
		return errors.Wrap(dbt.Error, "error creating exclusive partner record")
	}
	return nil
}

// GetExclusiveOffers gets a all exclusive partner records
func (r *UserRepository) GetExclusiveOffers(ctx context.Context) ([]entities.Exclusive, error) {
	var exclusives []entities.Exclusive
	dbt := r.DB.Raw(`SELECT * FROM exclusives WHERE partner_id IN (
        SELECT id FROM users where active = true
    )`).Scan(&exclusives)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.Wrap(dbt.Error, "error getting exclusive partner records")
	}
	return exclusives, nil
}

// GetExclusiveOffer gets the exclusive offer with the given partner ID
func (r *UserRepository) GetExclusiveOffer(ctx context.Context, partnerID uint) (*entities.Exclusive, error) {
	var exclusive entities.Exclusive
	dbt := r.DB.Where("partner_id = ?", partnerID).Find(&exclusive)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.Wrap(dbt.Error, "error getting exclusive partner record")
	}
	return &exclusive, nil
}

// DeleteExclusiveOffer deletes the given exclusive offer record
func (r *UserRepository) DeleteExclusiveOffer(ctx context.Context, exclusive *entities.Exclusive) (*entities.Exclusive, error) {
	dbt := r.DB.Delete(exclusive)
	if dbt.Error != nil {
		return nil, errors.Wrap(dbt.Error, "error deleting exclusive partner record")
	}
	return exclusive, nil
}

// CreateShare creates a share database record
func (r *UserRepository) CreateShare(ctx context.Context, share *entities.Share) error {
	var s entities.Share
	dbt := r.DB.Where("customer_id = ?", share.CustomerID).First(&s)
	if !dbt.RecordNotFound() {
		return errors.New("customer already shared")
	}
	dbt = r.DB.Create(&share)
	if dbt.Error != nil {
		return errors.Wrap(dbt.Error, "error creating share")
	}
	return nil
}

// GetSharesByCustomer returns 1 if there is a share from the given customer for the given partner in the given time range
func (r *UserRepository) GetSharesByCustomer(ctx context.Context, customerID uint, startDate time.Time, endDate time.Time) (int, error) {
	var share entities.Share
	dbt := r.DB.Where("customer_id = ? AND created_at > ? AND created_at < ?", customerID, startDate, endDate).First(&share)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return 0, nil
		}
		return 0, errors.Wrap(dbt.Error, "error getting share records")
	}
	return 1, nil
}

// SearchUsers searches the users database for the given search term. Search term can ba an email or phone number.
func (r *UserRepository) SearchUsers(ctx context.Context, searchTerm string) ([]entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var users []entities.User
	searchTerm = "%" + searchTerm + "%"
	dbt := r.DB.Where("email LIKE ? OR mobile LIKE ?", searchTerm, searchTerm).Find(&users)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			cancelFunc()
			return nil, nil
		}
		cancelFunc()
		return users, dbt.Error
	}
	cancelFunc()
	return users, nil
}

func (r UserRepository) GetSharablePartners(ctx context.Context) ([]entities.User, error) {
	var sharables []entities.User
	dbt := r.DB.Raw(`SELECT * FROM users WHERE id IN (SELECT partner_id FROM partner_profiles WHERE is_sharable = true)`).Scan(&sharables)
	if dbt.Error != nil {
		if dbt.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.Wrap(dbt.Error, "error getting share records")
	}
	return sharables, nil
}
