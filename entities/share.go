package entities

import (
	"github.com/jinzhu/gorm"
)

// Share represents a single share by a customer
type Share struct {
	gorm.Model      // Gorm DB model
	CustomerID uint `json:"customerID"`
}
