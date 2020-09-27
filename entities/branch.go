package entities

import (
	"github.com/jinzhu/gorm"
)

// Branch represents the branch belonging to a partner
type Branch struct {
	gorm.Model
	Country    string `gorm:"not null" json:"country"`
	CityID     uint   `json:"cityID"`
	City       City   `gorm:"not null" json:"city"`
	Address    string `gorm:"not null" json:"address"`
	Phone      string `gorm:"not null" json:"phone"`
	Mobile     string `json:"mobile"`
	OwnerID    uint   `gorm:"not null" json:"ownerID"`
	Owner      *User  `json:"owner"`
	CategoryID uint   `json:"categoryID"`
}
