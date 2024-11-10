package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type PaymentTransaction struct {
	ID            string  `validate:"required,uuid"`
	Amount        float64 `validate:"required,gte=0"`
	Currency      string  `validate:"required,oneof=USD EUR JPY"`
	CustomerID    string  `validate:"required,uuid"`
	Status        string  `validate:"required,oneof=pending completed failed"`
	PaymentMethod string  `validate:"required,oneof=credit debit paypal"`
	CreatedAt     time.Time
}

// NewPaymentTransaction creates a new PaymentTransaction if the input data is valid.
func NewPaymentTransaction(
	id string,
	amount float64,
	currency string,
	customerID string,
	method string,
	status string,
) (*PaymentTransaction, error) {
	pt := &PaymentTransaction{
		ID:            id,
		Amount:        amount,
		Currency:      currency,
		CustomerID:    customerID,
		Status:        status,
		PaymentMethod: method,
		CreatedAt:     time.Now(),
	}
	validate := validator.New()
	if err := validate.Struct(pt); err != nil {
		return nil, err
	}
	return pt, nil
}
