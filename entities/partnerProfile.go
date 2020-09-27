package entities

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// IPartnerProfile represents a partner profile interface
type IPartnerProfile interface {
	Approve(currentUser IUser) error
	SetDiscountValue(value float64)
	SetCategoryID(categoryID uint)
	SetMainBranchAddress(address string)
	SetPhone(phone string)
	SetCity(city *City)
	SetCountry(country string)
	SetBrandName(name string)
	SetLicenseURL(url string)
	SetLicenseKey(key string)
	SetOfferDescription(description string)
	AddPartnerPhoto(photo *PartnerPhoto)
	SetIsSharable(currentUser IUser, val bool) error
}

// PartnerProfile belongs to a user of type partner
type PartnerProfile struct {
	gorm.Model
	PartnerID         uint           `json:"partnerID,omitempty"` // foreign key
	Approved          bool           `gorm:"default:false" json:"approved"`
	DiscountValue     float64        `json:"discountValue"`
	CategoryID        uint           `json:"categroryID"`
	MainBranchAddress string         `json:"mainBranchAddress,omitempty"`
	Phone             string         `json:"phone,omitempty"`
	City              City           `json:"city,omitempty"`
	CityID            uint           `json:"cityID,omitempty"`
	Country           string         `json:"country,omitempty"`
	BrandName         string         `json:"brandName,omitempty"`
	LicenseURL        string         `json:"licenceURL,omitempty"`
	LicenseKey        string         `json:"licenceKey,omitempty"`
	OfferDiscription  string         `json:"offerDiscription"`
	PartnerPhotos     []PartnerPhoto `json:"photos,omitempty" gorm:"foreignkey:PartnerProfileID"`
	IsSharable        bool           `json:"isSharable" gorm:"default:false"`
}

// Approve sets the approved propety to true
func (profile *PartnerProfile) Approve(currentUser IUser) error {
	if !currentUser.IsAdmin() {
		return errors.New("only admins are allowed to approve users")
	}
	profile.Approved = true
	return nil
}

// SetDiscountValue sets the discountValue property
func (profile *PartnerProfile) SetDiscountValue(value float64) {
	profile.DiscountValue = value
}

// SetCategoryID sets the categoryID property
func (profile *PartnerProfile) SetCategoryID(categoryID uint) {
	profile.CategoryID = categoryID
}

// SetMainBranchAddress sets the main branch address property
func (profile *PartnerProfile) SetMainBranchAddress(address string) {
	profile.MainBranchAddress = address
}

// SetPhone sets the phone property
func (profile *PartnerProfile) SetPhone(phone string) {
	profile.Phone = phone
}

// SetCity sets the city property
func (profile *PartnerProfile) SetCity(city *City) {
	profile.City = *city
}

// SetCountry sets the country property
func (profile *PartnerProfile) SetCountry(country string) {
	profile.Country = country
}

// SetBrandName sets the brandName property
func (profile *PartnerProfile) SetBrandName(name string) {
	profile.BrandName = name
}

// SetLicenseURL sets the LicenseURL property
func (profile *PartnerProfile) SetLicenseURL(url string) {
	profile.LicenseURL = url
}

// SetLicenseKey sets the licenceKey property
func (profile *PartnerProfile) SetLicenseKey(key string) {
	profile.LicenseKey = key
}

// SetOfferDescription sets the offerDiscription property
func (profile *PartnerProfile) SetOfferDescription(description string) {
	profile.OfferDiscription = description
}

// AddPartnerPhoto adds a partner photo to the partner profile
func (profile *PartnerProfile) AddPartnerPhoto(photo *PartnerPhoto) {
	profile.PartnerPhotos = append(profile.PartnerPhotos, *photo)
}

// SetIsSharable sets the is sharable property of the partner profile
func (profile *PartnerProfile) SetIsSharable(currentUser IUser, val bool) error {
	if !currentUser.IsAdmin() {
		return errors.New("only admins are allowed to set a partner to be shared")
	}
	profile.IsSharable = val
	return nil
}
