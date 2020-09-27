package entities

import (
	"github.com/jinzhu/gorm"
	"time"
)

// OTP represents a one time password object
type OTP struct {
	gorm.Model
	UserID         uint
	HashedPassword []byte
	ExpiryDate     time.Time
}

// Expired returns true if tha expiry time of the password has passed
func (password *OTP) Expired() bool {
	return time.Now().After(password.ExpiryDate)
}

// ValidatePassword validates the given password agains the OTP hash
func (password *OTP) ValidatePassword(text string) bool {
	if password.Expired() {
		return false
	}
	return HashMatch(password.HashedPassword, text)
}
