package entities

import (
	"github.com/jinzhu/gorm"
)

// Category represents a category of branches
type Category struct {
	gorm.Model
	EnglishName string   `json:"englishName"`
	TurkishName string   `json:"turkishName"`
	Branches    []Branch `json:"-"`
}

// GetID is the getter for ID property
func (c *Category) GetID() uint {
	return c.ID
}

// AppendBranch adds a branch to the category
func (c *Category) AppendBranch(branch Branch) {
	c.Branches = append(c.Branches, branch)
}
