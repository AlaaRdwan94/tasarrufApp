package usecase

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/ahmedaabouzied/tasarruf/branch"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/offer"
	"github.com/ahmedaabouzied/tasarruf/subscription"
	"github.com/ahmedaabouzied/tasarruf/user"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
)

// OfferUsecase is an implementation of offer usecase interface
type OfferUsecase struct {
	offerRepo        offer.Repository
	hub              offer.Hub
	userRepo         user.Repository
	branchRepo       branch.Repository
	subscriptionRepo subscription.Repository
}

// CreateOfferUsecase returns an implementation of offer usecase interface
func CreateOfferUsecase(o offer.Repository, h offer.Hub, u user.Repository, b branch.Repository, s subscription.Repository) offer.Usecase {
	usecase := OfferUsecase{
		offerRepo:        o,
		hub:              h,
		userRepo:         u,
		branchRepo:       b,
		subscriptionRepo: s,
	}
	return &usecase
}

// ConsumeOffer represents a scan action
func (u *OfferUsecase) ConsumeOffer(ctx context.Context, customerID uint, partnerID uint, amount float64) (*entities.Offer, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	// Get current partner
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.getPartnerByID(ctx, currentUserID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if currentUser.ID != partnerID {
		err := errors.New("partner ID and current user ID do not match")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	// Get CustomerUser
	customer, err := u.getCustomerByID(ctx, customerID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if customer.Subscription == nil {
		err := errors.New("customer is not subscribed to any plan")
		log.Error(err)
		cancelFunc()
		return nil, err

	}
	if customer.Subscription.HasExpirPassed() || customer.Subscription.IsExpired() {
		err := errors.New("subscription has expired, please renew or upgrade your subscription")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	partner, err := u.getPartnerByID(ctx, partnerID)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	currentRemainingOffers, err := u.subscriptionRepo.GetCountOfOffersWithPartner(ctx, partner, customer.Subscription)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	customer.Subscription.RemainingOffers = currentRemainingOffers.CountOfOffers
	if !u.hub.HasUser(customer.ID) {
		err := errors.New("customer is not connected")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	net := amount - (currentUser.PartnerProfile.DiscountValue * amount / 100)
	offer := &entities.Offer{
		CustomerID:    customer.ID,
		PartnerID:     currentUser.ID,
		SubsriptionID: customer.Subscription.ID,
		Amount:        amount,
		Discount:      currentUser.PartnerProfile.DiscountValue,
		Total:         net,
		Customer:      customer,
		Partner:       currentUser,
	}
	err = currentUser.ConsumeOffer(customer, offer)
	if err != nil {
		err := errors.Wrap(err, "error consuming offer")
		log.Error(customer.Subscription.Expired, customer.Subscription.ID, customer.Subscription.ExpireDate)
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	// Save offer to DB
	offer, err = u.offerRepo.Create(ctx, offer)
	if err != nil {
		err := errors.Wrap(err, "repository error while creating offer")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	// set new remaining offers count
	err = u.subscriptionRepo.SetCountOfOffersWithPartner(ctx, partner, customer.Subscription, currentRemainingOffers.CountOfOffers-1)
	if err != nil {
		err := errors.Wrap(err, "repository error while decrementing remaining offers")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	// Send offer receipt to user
	err = u.hub.SendOfferToUser(customer.ID, offer)
	if err != nil {
		err := errors.Wrap(err, "WS error while sending notification to the customer")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return offer, nil
}

func (u *OfferUsecase) GetByCustomer(ctx context.Context, customerID uint) ([]entities.Offer, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	customer, err := u.getCustomerByID(ctx, customerID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting offers history")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	offers, err := u.offerRepo.GetByUser(ctx, customerID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting offers history")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	for i, offer := range offers {
		partner, err := u.getPartnerByID(ctx, offer.PartnerID)
		if err != nil {
			err := errors.Wrap(err, "repository error while getting offer partner")
			log.Error(err)
			cancelFunc()
			return nil, err
		}
		offers[i].Partner = partner
		offers[i].Customer = customer
	}
	cancelFunc()
	return offers, nil
}

func (u *OfferUsecase) GetMyOffersHistory(ctx context.Context, startDate time.Time, endDate time.Time) ([]entities.Offer, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	// Get current partner
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.userRepo.GetByID(ctx, currentUserID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting user")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	switch currentUser.AccountType {
	case "user":
		customer, err := u.getCustomerByID(ctx, currentUserID)
		if err != nil {
			err := errors.Wrap(err, "repository error while getting offers history")
			log.Error(err)
			cancelFunc()
			return nil, err
		}
		offers, err := u.offerRepo.GetByUser(ctx, currentUserID)
		if err != nil {
			err := errors.Wrap(err, "repository error while getting offers history")
			log.Error(err)
			cancelFunc()
			return nil, err
		}
		for i, offer := range offers {
			partner, err := u.getPartnerByID(ctx, offer.PartnerID)
			if err != nil {
				err := errors.Wrap(err, "repository error while getting offer partner")
				log.Error(err)
				cancelFunc()
				return nil, err
			}
			offers[i].Partner = partner
			offers[i].Customer = customer
		}
		cancelFunc()
		return offers, nil
	case "partner":
		partner, err := u.getPartnerByID(ctx, currentUserID)
		if err != nil {
			err := errors.Wrap(err, "repository error while getting offers history")
			log.Error(err)
			cancelFunc()
			return nil, err
		}
		offers, err := u.offerRepo.GetByPartner(ctx, currentUserID, startDate, endDate)
		if err != nil {
			err := errors.Wrap(err, "repository error while getting offers history")
			log.Error(err)
			cancelFunc()
			return nil, err
		}
		for i, offer := range offers {
			customer, err := u.getCustomerByID(ctx, offer.CustomerID)
			if err != nil {
				err := errors.Wrap(err, "repository error while getting offer partner")
				log.Error(err)
				cancelFunc()
				return nil, err
			}
			offers[i].Customer = customer
			offers[i].Partner = partner
		}
		cancelFunc()
		return offers, nil
	default:
		err := errors.New("repository error while getting offers history")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
}

// SendOffersStaticMail sends a static mail with the offers given by the provided date range
func (u *OfferUsecase) SendOffersStaticMail(ctx context.Context, startDate time.Time, endDate time.Time) error {
	offers, err := u.GetMyOffersHistory(ctx, startDate, endDate)
	if err != nil {
		return err
	}
	userID := ctx.Value(entities.UserIDKey).(uint)
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "error getting current user")
	}
	var records [][]string
	records = append(records, []string{"Date", "Customer Name", "Partner Name", "Amount", "Discount", "Total"})
	for _, offer := range offers {
		records = append(records, []string{offer.CreatedAt.Format("2 Jan 2006 15:04"), fmt.Sprintf("%s %s", offer.Customer.FirstName, offer.Customer.LastName), offer.Partner.PartnerProfile.BrandName, fmt.Sprintf("%.2f", offer.Amount), fmt.Sprintf("%.0f %%", offer.Discount), fmt.Sprintf("%.2f", offer.Total)})
	}
	buff := new(bytes.Buffer)
	w := csv.NewWriter(buff)
	err = w.WriteAll(records) // calls Flush internally
	if err != nil {
		return errors.Wrap(err, "failed to make csv file")
	}
	// create new *SGMailV3
	m := mail.NewV3Mail()
	from := mail.NewEmail("Tasarruf", "noreply@tasarruf.com")
	content := mail.NewContent("text/html", fmt.Sprintf("<h1>Offers Summary</h1> </br> <p> This is the summary of your offers on Tasarruf mobile application in the period from %s to %s\n", startDate.Format("2 Jan 2006"), endDate.Format("2 Jan 2006")))
	to := mail.NewEmail(fmt.Sprintf("%s %s", user.FirstName, user.LastName), user.Email)
	m.SetFrom(from)
	m.AddContent(content)
	// create new *Personalization
	personalization := mail.NewPersonalization()
	personalization.AddTos(to)
	personalization.Subject = "Tasarruf Summary"

	// add `personalization` to `m`
	m.AddPersonalizations(personalization)

	file := mail.NewAttachment()
	encoded := base64.StdEncoding.EncodeToString(buff.Bytes())
	file.SetContent(encoded)
	file.SetType("text/csv")
	file.SetFilename("summary.csv")
	file.SetDisposition("attachment")
	file.SetContentID("summary")
	m.AddAttachment(file)
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return nil
}

// GetOffer returns the offer with the given ID
func (u *OfferUsecase) GetOffer(ctx context.Context, offerID uint) (*entities.Offer, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	offer, err := u.offerRepo.GetByID(ctx, offerID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting offers history")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	customer, err := u.getCustomerByID(ctx, offer.CustomerID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting offer customer")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	offer.Customer = customer
	partner, err := u.getPartnerByID(ctx, offer.PartnerID)
	if err != nil {
		err := errors.Wrap(err, "repository error while getting offer partner")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	offer.Partner = partner
	cancelFunc()
	return offer, nil

}

// GetOffersCount returns the count of offers consumed
func (u *OfferUsecase) GetOffersCount(ctx context.Context) (int, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.userRepo.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return 0, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		err := errors.New("only admins are authorized to perform this task")
		cancelFunc()
		return 0, err
	}
	count, err := u.offerRepo.GetOffersCount(ctx)
	if err != nil {
		cancelFunc()
		return 0, errors.Wrap(err, "error getting count of customers")
	}
	cancelFunc()
	return count, nil
}

// GetAllOffers returns the count of offers consumed
func (u *OfferUsecase) GetAllOffers(ctx context.Context) ([]entities.Offer, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := u.userRepo.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		err := errors.New("only admins are authorized to perform this task")
		cancelFunc()
		return nil, err
	}
	offers, err := u.offerRepo.GetAllOffers(ctx)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting count of customers")
	}
	cancelFunc()
	return offers, nil
}

func (u *OfferUsecase) getCustomerByID(ctx context.Context, ID uint) (*entities.Customer, error) {
	user, err := u.userRepo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	otp, err := u.userRepo.GetOTPByUser(ctx, ID)
	if err != nil {
		return nil, err
	}
	subscription, err := u.subscriptionRepo.GetSubscriptionByUser(ctx, ID)
	if err != nil {
		return nil, err
	}
	if subscription != nil {
		plan, err := u.subscriptionRepo.GetPlanByID(ctx, subscription.PlanID)
		if err != nil {
			return nil, err
		}
		subscription.Plan = *plan
	}
	customer := entities.Customer{
		User:         *user,
		OTP:          otp,
		Subscription: subscription,
	}
	return &customer, nil
}

func (u *OfferUsecase) getPartnerByID(ctx context.Context, ID uint) (*entities.Partner, error) {
	user, err := u.userRepo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	otp, err := u.userRepo.GetOTPByUser(ctx, ID)
	if err != nil {
		return nil, err
	}
	partnerProfile, err := u.userRepo.GetPartnerProfileByUserID(ctx, ID)
	if err != nil {
		return nil, err
	}
	partner := entities.Partner{
		User:           *user,
		OTP:            otp,
		PartnerProfile: partnerProfile,
	}
	return &partner, nil
}
