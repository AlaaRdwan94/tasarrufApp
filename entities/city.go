package entities

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// ICity represents a city interface
type ICity interface {
	SetEnglishName(currentUser IUser, name string) error
	SetTurkishName(currentUser IUser, name string) error
}

// City holds city information.
// A city links a customer, partner and a branch.
type City struct {
	gorm.Model
	EnglishName string `json:"englishName"`
	TurkishName string `json:"turkishName"`
}

// SetEnglishName sets the english name property
func (c *City) SetEnglishName(currentUser IUser, name string) error {
	if !currentUser.IsAdmin() {
		return errors.New("not authorized")
	}
	c.EnglishName = name
	return nil
}

// SetTurkishName sets the turkish name property
func (c *City) SetTurkishName(currentUser IUser, name string) error {
	if !currentUser.IsAdmin() {
		return errors.New("not authorized")
	}
	c.TurkishName = name
	return nil
}
