package entities

import (
	"github.com/jinzhu/gorm"
	"time"
)

// ISubscription represents a subscription interface
type ISubscription interface {
	Expire()
	HasExpirPassed()
	IsExpired()
	GetRemainingOffers() uint
	HasRemainingOffers() bool
	AddRemainingOffers(amount uint)
	SubstractRemainingOffers(amount uint)
	SetPaymentID(paymentID string)
	SetPlan(plan *Plan)
	GetPlan() *Plan
}

// Subscription defines a user subscription to a plan
type Subscription struct {
	gorm.Model
	UserID              uint      `gorm:"not null" json:"userID"`
	PlanID              uint      `gorm:"not null" json:"planID"`
	ExpireDate          time.Time `gorm:"not null" json:"expireDate"`
	DelegationStartDate time.Time ` json:"-"`
	RemainingOffers     uint      `gorm:"not null" json:"remainingOffers"`
	Expired             bool      `gorm:"not null,default:fasle" json:"expired"`
	PaymentID           string    `json:"omit"`
	Plan                Plan      `json:"plan"`
}

// Expire the subscription
func (s *Subscription) Expire() {
	s.Expired = true
}

// HasExpirPassed returns true if the expire date has passed
func (s *Subscription) HasExpirPassed() bool {
	return time.Now().After(s.ExpireDate)
}

// IsExpired returns true if the expired field of the subscription is true
func (s *Subscription) IsExpired() bool {
	return s.Expired
}

// GetRemainingOffers returns the remaining offers
func (s *Subscription) GetRemainingOffers() uint {
	return s.RemainingOffers
}

// HasRemainingOffers returns true if the subscription remaining offers > 0
func (s *Subscription) HasRemainingOffers() bool {
	return s.RemainingOffers > 0
}

// AddRemainingOffers adds the given amount of remaining offers to the subscription
func (s *Subscription) AddRemainingOffers(amount uint) {
	s.RemainingOffers += amount
}

// SubstractRemainingOffers substracts the given amount of remaining offers from the subscription
func (s *Subscription) SubstractRemainingOffers(amount uint) {
	if s.HasRemainingOffers() {
		s.RemainingOffers -= amount
	}
}

// SetPaymentID sets the payment id property to the subscription
func (s *Subscription) SetPaymentID(ID string) {
	s.PaymentID = ID
}

// SetPlan sets the plan property of the subscription
func (s *Subscription) SetPlan(plan *Plan) {
	s.Plan = *plan
}

// GetPlan gets the plan property
func (s *Subscription) GetPlan() *Plan {
	return &s.Plan
}
