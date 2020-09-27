package entities

import (
	"testing"
)

func TestSetEnglishName(t *testing.T) {
	t.Run("ByNotAdmin", func(t *testing.T) {
		user := User{
			AccountType: "user",
		}
		city := City{}
		err := city.SetEnglishName(&user, "updated")
		if err == nil {
			t.Fail()
		}
	})
	t.Run("ByAdmin", func(t *testing.T) {
		user := User{
			AccountType: "admin",
		}
		city := City{}
		err := city.SetEnglishName(&user, "updated")
		if err != nil {
			t.Error(err)
		}
	})
}

func TestSetTurkishName(t *testing.T) {
	t.Run("ByNotAdmin", func(t *testing.T) {
		user := User{
			AccountType: "user",
		}
		city := City{}
		err := city.SetTurkishName(&user, "updated")
		if err == nil {
			t.Fail()
		}
	})
	t.Run("ByAdmin", func(t *testing.T) {
		user := User{
			AccountType: "admin",
		}
		city := City{}
		err := city.SetTurkishName(&user, "updated")
		if err != nil {
			t.Error(err)
		}
	})
}
