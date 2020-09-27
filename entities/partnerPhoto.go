package entities

import (
	"github.com/jinzhu/gorm"
)

// PartnerPhoto represents a photo belonging to a partner profile
type PartnerPhoto struct {
	gorm.Model
	PartnerProfileID uint   `json:"partnerProfileID"`
	PhotoURL         string `json:"photoURL"`
	PhotoKey         string `json:"photoKey"`
}
