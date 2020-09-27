package entities

import (
	"github.com/jinzhu/gorm"
)

// ICustomer represents a customer interface
type CustomerPartnerOffersCount struct {
	gorm.Model
	CustomerID     uint
	PartnerID      uint
	SubscriptionID uint
	CountOfOffers  uint
}
