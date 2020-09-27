package entities

import (
	"github.com/pkg/errors"
)

// IPartner represents a partner interface
type IPartner interface {
	SetPartnerProfile(profile *PartnerProfile)
	GetPartnerProfile() *PartnerProfile
	Approve(currentUser IUser) error
	SetDiscountValue(currentUser IUser, value float64) error
	SetCategoryID(currentUser IUser, categroryID uint) error
	SetMainBranchAddress(currentUser IUser, address string) error
	SetPhone(currentUser IUser, phone string) error
	SetCity(currentUser IUser, city *City) error
	SetCountry(currentUser IUser, country string) error
	SetBrandName(currentUser IUser, name string) error
	SetLicenseURL(currentUser IUser, URL string) error
	SetLicenseKey(currentUser IUser, key string) error
	SetOfferDescription(currentUser IUser, description string) error
	AddPartnerPhoto(currentUser IUser, photo *PartnerPhoto) error
	ConsumeOffer(customer ICustomer, offer *Offer) error
}

// Partner is a user of type partner
type Partner struct {
	User                           // embeded
	PartnerProfile *PartnerProfile `json:"partnerProfile"`
	OTP            *OTP
	Offers         []Offer
	Reviews        []Review
}

// SetPartnerProfile sets the partner profile property
func (partner *Partner) SetPartnerProfile(profile *PartnerProfile) {
	partner.PartnerProfile = profile
}

// GetPartnerProfile gets the partner profile property
func (partner *Partner) GetPartnerProfile() *PartnerProfile {
	return partner.PartnerProfile
}

// Approve sets the partner profile approved property
func (partner *Partner) Approve(currentUser IUser) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}
	err = partner.PartnerProfile.Approve(currentUser)
	if err != nil {
		return err
	}
	return nil
}

// SetDiscountValue sets the discountValue property
func (partner *Partner) SetDiscountValue(currentUser IUser, value float64) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}
	partner.PartnerProfile.SetDiscountValue(value)
	return nil
}

// SetCategoryID sets the categroryID value
func (partner *Partner) SetCategoryID(currentUser IUser, categroryID uint) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}
	partner.PartnerProfile.SetCategoryID(categroryID)
	return nil
}

// SetMainBranchAddress sets the mainBranchAddress property
func (partner *Partner) SetMainBranchAddress(currentUser IUser, address string) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}
	partner.PartnerProfile.SetMainBranchAddress(address)
	return nil
}

// SetPhone sets the phone property
func (partner *Partner) SetPhone(currentUser IUser, phone string) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}
	partner.PartnerProfile.SetPhone(phone)
	return nil
}

// SetCountry sets the user country
func (partner *Partner) SetCountry(currentUser IUser, newCountry string) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}
	partner.User.Country = newCountry
	partner.PartnerProfile.SetCountry(newCountry)
	return nil
}

// SetCity sets the user city
func (partner *Partner) SetCity(currentUser IUser, newCity *City) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}
	partner.User.City = *newCity
	partner.PartnerProfile.SetCity(newCity)
	return nil
}

// SetBrandName sets the brand property
func (partner *Partner) SetBrandName(currentUser IUser, name string) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}
	partner.PartnerProfile.SetBrandName(name)
	return nil
}

// SetLicenseURL sets the licenceURL property
func (partner *Partner) SetLicenseURL(currentUser IUser, URL string) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}

	partner.PartnerProfile.SetLicenseURL(URL)
	return nil
}

// SetLicenseKey sets the licenseKey property
func (partner *Partner) SetLicenseKey(currentUser IUser, key string) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}
	partner.PartnerProfile.SetLicenseKey(key)
	return nil
}

// SetOfferDescription sets the offerDiscription property
func (partner *Partner) SetOfferDescription(currentUser IUser, description string) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}
	partner.PartnerProfile.SetOfferDescription(description)
	return nil
}

// AddPartnerPhoto adds a partner photo
func (partner *Partner) AddPartnerPhoto(currentUser IUser, photo *PartnerPhoto) error {
	err := verifyAuthorization(partner, currentUser)
	if err != nil {
		return err
	}
	partner.PartnerProfile.AddPartnerPhoto(photo)
	return nil
}

// ConsumeOffer consumes an offer from the given customer
func (partner *Partner) ConsumeOffer(customer ICustomer, offer *Offer) error {
	if partner.PartnerProfile == nil {
		return errors.New("error getting partner profile")
	}
	if !partner.PartnerProfile.Approved {
		return errors.New("partner not approved")
	}
	return customer.AddOffer(offer)
}

// ToggleIsSharable toggles the value of is sharable of the customer profile
func (partner *Partner) ToggleIsSharable(currentUser IUser) error {
	if partner.PartnerProfile == nil {
		return errors.New("error getting partner profile")
	}
	if !partner.PartnerProfile.Approved {
		return errors.New("partner is not approved")
	}
	return partner.PartnerProfile.SetIsSharable(currentUser, !partner.PartnerProfile.IsSharable)
}
