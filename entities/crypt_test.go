package entities

import "testing"

func TestEncryptPassword(t *testing.T) {
	p := "1234566"
	_, err := EncryptPassword(p)
	if err != nil {
		t.Fail()
	}
}

func TestVerifyPassword(t *testing.T) {
	p := "123456"
	hp, err := EncryptPassword(p)
	if err != nil {
		t.Fail()
	}
	mockUser := User{
		HashedPassword: hp,
	}
	t.Run("VerifyTruePassword", func(t *testing.T) {
		v := mockUser.VerifyPassword(p)
		if v == false {
			t.Fail()
		}
	})
	t.Run("VerifyFalsePassword", func(t *testing.T) {
		v := mockUser.VerifyPassword("654321")
		if v == true {
			t.Fail()
		}
	})
}
