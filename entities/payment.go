package entities

import (
	"github.com/ahmedaabouzied/iyzipay-go/iyzipay"
)

// PaymentDetails represents a transaction
type PaymentDetails struct {
	User     *User
	Plan     *Plan
	Card     *iyzipay.PaymentCard
	IDNumber string
}
