package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ahmedaabouzied/iyzipay-go/iyzipay"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/pkg/errors"
	"os"
)

type transaction struct {
	user        *entities.User
	plan        *entities.Plan
	card        *iyzipay.PaymentCard
	idNumber    string
	iyzipayRoot iyzipay.Options
}

// CreateTransaction returns a new iyzipay transaction.
// An iyzipay transaction implements the payment interface.
func CreateTransaction(details *entities.PaymentDetails) Payment {
	apiKey := os.Getenv("PAYMENT_API_KEY")
	secretKey := os.Getenv("PAYMENT_API_SECRET")
	baseURL := os.Getenv("PAYMENT_BASE_URL")

	options := iyzipay.Options{}
	options.New(apiKey, secretKey, baseURL)
	transaction := transaction{
		iyzipayRoot: options,
		user:        details.User,
		plan:        details.Plan,
		idNumber:    details.IDNumber,
		card:        details.Card,
	}
	return &transaction
}

func (t *transaction) Submit(ctx context.Context) (string, error) {

	paymentCard := iyzipay.PaymentCard{
		CardHolderName: t.user.FirstName,
		CardNumber:     t.card.CardNumber,
		ExpireMonth:    t.card.ExpireMonth,
		ExpireYear:     t.card.ExpireYear,
		Cvc:            t.card.Cvc,
	}

	address := iyzipay.Address{
		ContactName: fmt.Sprintf("%s %s", t.user.FirstName, t.user.LastName),
		Address:     fmt.Sprintf("%s %s", t.user.City.EnglishName, t.user.Country),
		City:        t.user.City.EnglishName,
		Country:     t.user.Country,
		ZipCode:     "000000",
	}

	buyer := iyzipay.Buyer{
		Id:                  fmt.Sprintf("%d", t.user.ID),
		Name:                t.user.FirstName,
		Surname:             t.user.LastName,
		Email:               t.user.Email,
		GsmNumber:           t.user.Mobile,
		IdentityNumber:      t.idNumber,
		RegistrationAddress: address.Address,
		City:                address.City,
		Country:             address.Country,
		ZipCode:             address.ZipCode,
	}

	basketItems := []iyzipay.BasketItem{
		iyzipay.BasketItem{
			Id:        fmt.Sprintf("%d", t.plan.ID),
			Name:      t.plan.EnglishName,
			Category1: "subscription plan",
			ItemType:  "VIRTUAL",
			Price:     fmt.Sprintf("%f", t.plan.Price),
		},
	}

	request := iyzipay.CreatePaymentRequest{
		Locale:          "en",
		ConversationId:  "123456789",
		Price:           fmt.Sprintf("%f", t.plan.Price),
		PaidPrice:       fmt.Sprintf("%f", t.plan.Price),
		BasketId:        fmt.Sprintf("%d", t.plan.ID),
		PaymentGroup:    "LISTING",
		PaymentCard:     paymentCard,
		Currency:        "TRY",
		Buyer:           buyer,
		ShippingAddress: address,
		BillingAddress:  address,
		BasketItems:     basketItems,
	}

	paymentResponse := iyzipay.Payment{}.Create(request, t.iyzipayRoot)
	var resp map[string]interface{}
	err := json.Unmarshal([]byte(paymentResponse), &resp)
	if err != nil {
		return "", errors.Wrap(err, "error processing payment: error processing response from Iyzipay")
	}
	if status, ok := resp["status"]; ok {
		if status.(string) == "success" {
			return "success", nil
		}
		if status.(string) == "failure" {
			return resp["errorMessage"].(string), errors.New(resp["errorMessage"].(string))
		}
	}
	return "error processing payment", errors.New("error processing payment")
}

func (t *transaction) Cancel(ctx context.Context) (string, error) {
	return "", nil
}
