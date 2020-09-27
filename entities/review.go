package entities

import (
	"github.com/jinzhu/gorm"
)

// Review represents a review object from a customer to a partner.
// Belongs to customer and a partner.
type Review struct {
	gorm.Model        // Gorm DB model
	CustomerID uint   `json:"customerID"`
	PartnerID  uint   `json:"partnerID"`
	Stars      int    `json:"stars"`
	Content    string `json:"content"`
	Customer   User   `json:"customer,omitempty"`
	Partner    User   `json:"partner,omitempty"`
}
