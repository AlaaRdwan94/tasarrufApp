package entities

import (
	"github.com/jinzhu/gorm"
)

// Plan represents the plan a user subscribes on
type Plan struct {
	gorm.Model
	EnglishName        string  `gorm:"not null" json:"englishName"`
	EnglishDescription string  `json:"engishDescription"`
	TurkishDescription string  `json:"turkishDescription"`
	TurkishName        string  `gorm:"not null" json:"trukishName"`
	Price              float64 `gorm:"not null" json:"price"`
	CountOfOffers      uint    `gorm:"not null" json:"countOfOffers"`
	Image              string  `json:"image"`
	Rank               uint    `json:"rank"`      // rank determines the order of which the plan is displayed on clients
	IsDefault          bool    `json:"isDefault"` // default plan is the plan for users with expired plans and new users
}
