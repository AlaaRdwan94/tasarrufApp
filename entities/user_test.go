package entities

import (
	"testing"
	"time"
)

func TestValidateTruePasswordWithoutOTP(t *testing.T) {
	hash, err := EncryptPassword("12345678")
	if err != nil {
		t.Error(err)
		return
	}
	u := User{
		HashedPassword: hash,
	}
	err = u.ValidatePassword("12345678")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestValidateFalsePasswordWithoutOTP(t *testing.T) {
	hash, err := EncryptPassword("12345678")
	if err != nil {
		t.Error(err)
	}
	u := User{
		HashedPassword: hash,
	}
	err = u.ValidatePassword("12345677")
	if err == nil {
		t.Fail()
		return
	}
}

func TestValidateTruePasswordWithOTP(t *testing.T) {
	hash, err := EncryptPassword("12345678")
	if err != nil {
		t.Error(err)
	}
	u := User{
		HashedPassword: hash,
		OTP: &OTP{
			HashedPassword: []byte{},
		},
	}
	err = u.ValidatePassword("12345678")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestValidateFalsePasswordWithOTP(t *testing.T) {
	hash, err := EncryptPassword("12345678")
	if err != nil {
		t.Error(err)
	}
	u := User{
		HashedPassword: hash,
		OTP: &OTP{
			HashedPassword: []byte{},
		},
	}
	err = u.ValidatePassword("12345677")
	if err == nil {
		t.Fail()
		return
	}
}

func TestValidateTrueNotExpiredOTPPassword(t *testing.T) {
	hash, err := EncryptPassword("12345678")
	if err != nil {
		t.Error(err)
		return
	}
	otpHash, err := EncryptPassword("otpPass")
	if err != nil {
		t.Error(err)
		return
	}
	u := User{
		HashedPassword: hash,
		OTP: &OTP{
			HashedPassword: otpHash,
			ExpiryDate:     time.Now().Add(5 * time.Second),
		},
	}
	err = u.ValidatePassword("otpPass")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestValidateTrueExpiredOTPPassword(t *testing.T) {
	hash, err := EncryptPassword("12345678")
	if err != nil {
		t.Error(err)
		return
	}
	otpHash, err := EncryptPassword("otpPass")
	if err != nil {
		t.Error(err)
		return
	}
	u := User{
		HashedPassword: hash,
		OTP: &OTP{
			HashedPassword: otpHash,
			ExpiryDate:     time.Now(),
		},
	}
	err = u.ValidatePassword("otpPass")
	if err == nil {
		t.Fail()
		return
	}
}

func TestValidateFalseOTPPassword(t *testing.T) {
	hash, err := EncryptPassword("12345678")
	if err != nil {
		t.Error(err)
		return
	}
	otpHahs, err := EncryptPassword("otpPass")
	if err != nil {
		t.Error(err)
		return
	}
	u := User{
		HashedPassword: hash,
		OTP: &OTP{
			HashedPassword: otpHahs,
		},
	}
	err = u.ValidatePassword("false")
	if err == nil {
		t.Fail()
		return
	}
}

func TestSetDiffirentPassword(t *testing.T) {
	hash, err := EncryptPassword("12345678")
	if err != nil {
		t.Error(err)
	}
	u := User{
		HashedPassword: hash,
	}
	err = u.SetHashedPassword("12345677", &u)
	if err != nil {
		t.Error(err)
	}
}

func TestSetSamePassword(t *testing.T) {
	hash, err := EncryptPassword("12345678")
	if err != nil {
		t.Error(err)
	}
	u := User{
		HashedPassword: hash,
	}
	err = u.SetHashedPassword("12345678", &u)
	if err == nil {
		t.Fail()
	}
}

func TestToggleActiveByAdmin(t *testing.T) {
	user := User{Active: true}
	admin := User{AccountType: "admin"}
	err := user.ToggleActive(&admin)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if user.Active {
		t.Fail()
	}
}

func TestToggleActiveByNotAdmin(t *testing.T) {
	user := User{Active: true}
	admin := User{AccountType: "partner"}
	err := user.ToggleActive(&admin)
	if err == nil {
		t.Fail()
	}
}
