package entities

import (
	"github.com/jinzhu/gorm"
)

// PlanCategory represents a category plan many to many relationship
type PlanCategory struct {
	gorm.Model
	PlanID     uint `gorm:"not null"`
	CategoryID uint `gorm:"not null"`
}
