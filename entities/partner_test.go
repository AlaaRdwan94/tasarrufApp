package entities

import (
	"testing"
	"time"
)

func TestConsumeOfferWithExpiredSubscription(t *testing.T) {
	customer := Customer{
		User: User{},
		Subscription: &Subscription{
			Expired:         true,
			RemainingOffers: 1,
			ExpireDate:      time.Now().Add(60 * time.Second),
		},
	}
	partner := &Partner{
		PartnerProfile: &PartnerProfile{},
	}

	err := partner.ConsumeOffer(&customer, &Offer{})
	if err == nil {
		t.Fail()
	}
}

func TestConsumeOfferWithPassedExpireDate(t *testing.T) {
	customer := Customer{
		User: User{},
		Subscription: &Subscription{
			Expired:         false,
			RemainingOffers: 1,
			ExpireDate:      time.Now().Add(-60 * time.Second),
		},
	}
	partner := &Partner{
		PartnerProfile: &PartnerProfile{},
	}
	err := partner.ConsumeOffer(&customer, &Offer{})
	if err == nil {
		t.Fail()
	}
}

func TestConsumeOfferWithNoRemainingOffers(t *testing.T) {
	customer := Customer{
		User: User{},
		Subscription: &Subscription{
			Expired:         false,
			RemainingOffers: 0,
			ExpireDate:      time.Now().Add(60 * time.Second),
		},
	}
	partner := &Partner{
		PartnerProfile: &PartnerProfile{},
	}
	err := partner.ConsumeOffer(&customer, &Offer{})
	if err == nil {
		t.Fail()
	}
}

func TestDecrementRemainingOfferOnConsume(t *testing.T) {
	customer := Customer{
		User: User{
			Verified: true,
		},
		Subscription: &Subscription{
			Expired:         false,
			RemainingOffers: 1,
			ExpireDate:      time.Now().Add(60 * time.Second),
		},
	}
	partner := &Partner{
		PartnerProfile: &PartnerProfile{
			Approved: true,
		},
	}
	err := partner.ConsumeOffer(&customer, &Offer{})
	if err != nil {
		t.Error(err)
	}
	if customer.Subscription.RemainingOffers != 0 {
		t.Fail()
	}
}

func TestConsumeOfferWithNotApprovedPartner(t *testing.T) {
	customer := Customer{
		User: User{},
		Subscription: &Subscription{
			Expired:         false,
			RemainingOffers: 1,
			ExpireDate:      time.Now().Add(60 * time.Second),
		},
	}
	partner := &Partner{
		PartnerProfile: &PartnerProfile{
			Approved: false,
		},
	}
	err := partner.ConsumeOffer(&customer, &Offer{})
	if err == nil {
		t.Fail()
	}
}

func TestConsumeOfferWithNotVerifiedCustomer(t *testing.T) {
	customer := Customer{
		User: User{
			Verified: false,
		},
		Subscription: &Subscription{
			Expired:         false,
			RemainingOffers: 1,
			ExpireDate:      time.Now().Add(60 * time.Second),
		},
	}
	partner := &Partner{
		PartnerProfile: &PartnerProfile{
			Approved: true,
		},
	}
	err := partner.ConsumeOffer(&customer, &Offer{})
	if err == nil {
		t.Fail()
	}
}

func TestToggleIsSharableWithNotAdminUser(t *testing.T) {
	notAdmin := User{
		AccountType: "user",
	}
	partner := &Partner{
		PartnerProfile: &PartnerProfile{
			Approved: false,
		},
	}
	err := partner.ToggleIsSharable(&notAdmin)
	if err == nil {
		t.Fail()
	}
}

func TestToggleIsSharableWithAdminUser(t *testing.T) {
	notAdmin := User{
		AccountType: "admin",
	}
	partner := &Partner{
		PartnerProfile: &PartnerProfile{
			Approved: true,
		},
	}
	err := partner.ToggleIsSharable(&notAdmin)
	if err != nil {
		t.Error(err)
	}
}

func TestToggleIsSharableForNotApprovedPartner(t *testing.T) {
	notAdmin := User{
		AccountType: "admin",
	}
	partner := &Partner{
		PartnerProfile: &PartnerProfile{
			Approved: false,
		},
	}
	err := partner.ToggleIsSharable(&notAdmin)
	if err == nil {
		t.Fail()
	}
}
