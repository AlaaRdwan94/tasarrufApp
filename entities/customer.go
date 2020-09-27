package entities

import (
	"github.com/pkg/errors"
	"time"
)

// ICustomer represents a customer interface
type ICustomer interface {
	SetDateOfBirth(newDob time.Time, currentUser IUser) error
	GetDateOfBirth() time.Time
	SetCity(currentUser IUser, newCity *City) error
	SetCountry(currentUser IUser, newCountry string) error
	SubstractRemainingOffers(amount uint) error
	AddRemainingOffers(amount uint)
	AddOffer(offer *Offer) error
}

// Customer represents a user of customer type
type Customer struct {
	User         // embed
	OTP          *OTP
	Subscription *Subscription
	Reviews      []Review
}

// SetDateOfBirth sets the customer date of birth property
func (c *Customer) SetDateOfBirth(newDob time.Time, currentUser IUser) error {
	err := verifyAuthorization(c, currentUser)
	if err != nil {
		return errors.Wrap(err, "Not Authorized to set date of birth")
	}
	c.DateOfBirth = newDob
	return nil
}

// GetDateOfBirth returns the customer date of birth
func (c *Customer) GetDateOfBirth() time.Time {
	return c.User.DateOfBirth
}

// SetCity sets the city property
func (c *Customer) SetCity(currentUser IUser, newCity *City) error {
	err := verifyAuthorization(c, currentUser)
	if err != nil {
		return err
	}
	c.User.City = *newCity
	return nil
}

// SetCountry sets the country property
func (c *Customer) SetCountry(currentUser IUser, newCountry string) error {
	err := verifyAuthorization(c, currentUser)
	if err != nil {
		return err
	}
	c.User.Country = newCountry
	return nil
}

// AddRemainingOffers adds the amount to the remainingOffers
func (c *Customer) AddRemainingOffers(amount uint) {
	c.Subscription.AddRemainingOffers(amount)
}

// ExpireSubscription sets the subscription to be expired
func (c *Customer) ExpireSubscription() error {
	if c.Subscription == nil {
		return errors.New("user is not subscriped to any plan")
	}
	c.Subscription.Expire()
	return nil
}

// SubstractRemainingOffers substracts the amount from remaining offers
func (c *Customer) SubstractRemainingOffers(amount uint) error {
	if c.Subscription == nil {
		return errors.New("not subscriped to any plan")
	}
	if c.Subscription.HasExpirPassed() {
		return errors.New("expired subscription")
	}
	if c.Subscription.IsExpired() {
		return errors.New("expired subscription")
	}
	if !c.Subscription.HasRemainingOffers() {
		return errors.New("customer doesn't have remaining offers left")
	}
	c.Subscription.SubstractRemainingOffers(amount)
	return nil
}

// AddOffer adds an offer to the customer offers array property
func (c *Customer) AddOffer(offer *Offer) error {
	if !c.User.Verified {
		return errors.New("customer not verified")
	}
	err := c.SubstractRemainingOffers(1)
	if err != nil {
		return err
	}
	return nil
}
