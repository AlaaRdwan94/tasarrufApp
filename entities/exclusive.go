package entities

import (
	"github.com/jinzhu/gorm"
)

// Exclusive represnts an exclusive offer record
type Exclusive struct {
	gorm.Model
	PartnerID uint `json:"partnerID"`
}
