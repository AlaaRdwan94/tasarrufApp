package entities

import (
	"github.com/jinzhu/gorm"
)

// SupportInfo represents support contact information table
type SupportInfo struct {
	gorm.Model
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}
