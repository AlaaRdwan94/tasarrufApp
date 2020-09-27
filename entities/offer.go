package entities

import (
	"github.com/jinzhu/gorm"
)

// Offer represents an offer DB model
type Offer struct {
	gorm.Model
	CustomerID    uint      `json:"customerID,omitempty"`
	PartnerID     uint      `json:"partnerID,omitempty"`
	SubsriptionID uint      `json:"subscriptionID"`
	Amount        float64   `json:"amount,omitempty"`
	Discount      float64   `json:"discount,omitempty"`
	Total         float64   `json:"total,omitempty"`
	Customer      *Customer `json:"customer,omitempty" gorm:"-"`
	Partner       *Partner  `json:"partner,omitempty" gorm:"-"`
}
