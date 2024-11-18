package model

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

type PaymentTransaction struct {
	ID            string  `validate:"required,uuid"`
	Amount        float64 `validate:"required,gte=0"`
	Currency      string  `validate:"required,oneof=USD EUR JPY"`
	CustomerID    *string `validate:"omitempty,uuid"`
	Status        string  `validate:"required,oneof=pending completed failed"`
	PaymentMethod string  `validate:"required,oneof=credit debit paypal"`
	CreatedAt     time.Time
}

// NewPaymentTransaction creates a new PaymentTransaction if the input data is valid.
func NewPaymentTransaction(
	id string,
	amount float64,
	currency string,
	method string,
) (*PaymentTransaction, error) {
	pt := &PaymentTransaction{
		ID:            id,
		Amount:        amount,
		Currency:      currency,
		Status:        "pending",
		PaymentMethod: method,
		CreatedAt:     time.Now(),
	}
	validate := validator.New()
	if err := validate.Struct(pt); err != nil {
		return nil, err
	}
	return pt, nil
}

// BindCustomerToTransaction binds a customer ID to a transaction
func (pt *PaymentTransaction) BindCustomerToTransaction(cID string) {
	if pt.CustomerID == nil {
		pt.CustomerID = &cID
		return
	}

	log.Fatalf("Customer ID already bound to transaction: %s", pt.CustomerID)
}
