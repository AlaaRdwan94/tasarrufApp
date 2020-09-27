// Copyright 2019 NOVA Solutions Co. All Rights Reserved.
//

package entities

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

// Key type is not exported to not collide with other context keys defined in other packages
type key string

// UserIDKey is the context key for the user ID.
const UserIDKey key = "ID"

// IUser represents a User interfce
type IUser interface {
	SetEmail(newEmail string, currentUser IUser) error
	SetMobile(newMobile string, currentUser IUser) error
	SetHashedPassword(newPassword string, currentUser IUser) error
	SetHashedVerificationCode(newCode []byte, currentUser IUser) error
	SetProfileImageURL(newProileImageURL string, currentUser IUser) error
	SetProfileImageKey(newProfileImageKey string, currentUser IUser) error
	SetVerified(currentUser IUser) error
	SetOTP(otp *OTP)
	GetOTP() *OTP
	GetEmail() string
	GetMobile() string
	GetCountry() string
	GetCity() *City
	GetHashedPassword() []byte
	GetHashedVerificationCode() []byte
	GetProfileImageURL() string
	GetProfileImageKey() string
	IsVerified() bool
	IsAdmin() bool
	GetID() uint
	GetAccountType() string
	ToggleActive(admin IUser) error
}

// User represents the user of the app.  Can be admin or user.
type User struct {
	gorm.Model
	Email                  string         `gorm:"unique;not null" json:"email"`
	Mobile                 string         `gorm:"unique;not null" json:"mobile,omitempty"`
	FirstName              string         `json:"firstName"`
	LastName               string         `json:"lastName"`
	AccountType            string         `json:"accountType"`
	Country                string         `json:"country"`
	CityID                 uint           `json:"cityID"`
	DateOfBirth            time.Time      `json:"dateOfBirth"`
	HashedPassword         []byte         `json:"-"`
	HashedVerificationCode []byte         `json:"-"`
	ProfileImageURL        string         `json:"profileImageURL"`
	ProfileImageKey        string         `json:"profileImageKey"`
	PartnerProfile         PartnerProfile `gorm:"foreignkey:PartnerProfileID" json:"partnerProfile,omitempty"`
	Verified               bool           `gorm:"default:false" json:"verified"`
	OTP                    *OTP           `json:"-"`
	City                   City           `json:"city" gorm:"-"`
	Active                 bool           `json:"active" gorm:"default:true;not null"`
}

// IsPartner returns true if the user account type is partner
func (user *User) IsPartner() bool {
	return user.AccountType == "partner"
}

// IsAdmin return true if the user is admin
func (user *User) IsAdmin() bool {
	return user.AccountType == "admin"
}

// GetID gets the user ID
func (user *User) GetID() uint {
	return user.ID
}

// GetAccountType return the account type of the user
func (user *User) GetAccountType() string {
	return user.AccountType
}

// SetEmail sets the user Email
func (user *User) SetEmail(newEmail string, currentUser IUser) error {
	err := verifyAuthorization(user, currentUser)
	if err != nil {
		return err
	}
	user.Email = newEmail
	return nil
}

// SetMobile sets the user mobile
func (user *User) SetMobile(newMobile string, currentUser IUser) error {
	err := verifyAuthorization(user, currentUser)
	if err != nil {
		return err
	}
	user.Mobile = newMobile
	return nil
}

// SetHahshedPassword sets the user hashed password property
func (user *User) SetHashedPassword(newPassword string, currentUser IUser) error {
	err := verifyAuthorization(user, currentUser)
	if err != nil {
		return err
	}
	if HashMatch(user.HashedPassword, newPassword) {
		return errors.New("new password cannot be the same as old password")
	}
	hash, err := EncryptPassword(newPassword)
	if err != nil {
		return errors.Wrap(err, "error generating password hash")
	}
	user.HashedPassword = hash
	return nil
}

// SetHashedVerificationCode sets the user hashed verification code property
func (user *User) SetHashedVerificationCode(newCode []byte, currentUser IUser) error {
	err := verifyAuthorization(user, currentUser)
	if err != nil {
		return err
	}
	user.HashedVerificationCode = newCode
	return nil
}

// SetProfileImageURL sets the user profile image URL
func (user *User) SetProfileImageURL(newProileImageURL string, currentUser IUser) error {
	err := verifyAuthorization(user, currentUser)
	if err != nil {
		return err
	}
	user.ProfileImageURL = newProileImageURL
	return nil
}

// SetProfileImageKey sets the user profile image key
func (user *User) SetProfileImageKey(newProfileImageKey string, currentUser IUser) error {
	err := verifyAuthorization(user, currentUser)
	if err != nil {
		return err
	}
	user.ProfileImageKey = newProfileImageKey
	return nil
}

// SetVerified sets the user verified property
func (user *User) SetVerified(currentUser IUser) error {
	err := verifyAuthorization(user, currentUser)
	if err != nil {
		return err
	}
	user.Verified = true
	return nil
}

// GetEmail returns the user email
func (user *User) GetEmail() string {

	return user.Email
}

// GetMobile returns the user mobile
func (user *User) GetMobile() string {
	return user.Mobile
}

// GetCountry returns the user country
func (user *User) GetCountry() string {
	return user.Country
}

// GetCity returns the user city
func (user *User) GetCity() *City {
	return &user.City
}

// GetHashedPassword returns the user hashed password
func (user *User) GetHashedPassword() []byte {
	return user.HashedPassword
}

// GetHashedVerificationCode retuns the user hashed verification code
func (user *User) GetHashedVerificationCode() []byte {
	return user.HashedVerificationCode
}

// GetProfileImageURL returns the user profile image URL
func (user *User) GetProfileImageURL() string {
	return user.ProfileImageURL
}

// GetProfileImageKey returns the user profile image key
func (user *User) GetProfileImageKey() string {
	return user.ProfileImageKey
}

// IsVerified returns the user verified property
func (user *User) IsVerified() bool {
	return user.Verified
}

// GetOTP gets the user one time password object
func (user *User) GetOTP() *OTP {
	return user.OTP
}

// SetOTP sets the user one time password
func (user *User) SetOTP(otp *OTP) {
	user.OTP = otp
}

// ToggleActive toggles user active propery. Returns error if provided user is not admin.
func (user *User) ToggleActive(admin IUser) error {
	if !admin.IsAdmin() {
		return errors.New("not authorized to toggle user active property")
	}
	user.Active = !user.Active
	return nil
}

// ValidatePassword validates the password against the password and one time passwords
func (user *User) ValidatePassword(password string) error {
	if HashMatch(user.GetHashedPassword(), password) {
		return nil
	}
	otp := user.GetOTP()
	if otp == nil {
		return errors.New("invalid password")
	}
	if otp.ValidatePassword(password) {
		return nil
	}
	return errors.New("invalid password")
}

func verifyAuthorization(c1 IUser, c2 IUser) error {
	if c1.GetID() == c2.GetID() || c2.IsAdmin() {
		return nil
	}
	return errors.New("Not Authroized")
}
