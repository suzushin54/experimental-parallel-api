package gateway

import "context"

type PaymentGateway interface {
	ProcessPayment(ctx context.Context, userID string, amount float64) (bool, error)
}

type paymentGateway struct {
	amount   float64
	currency string
	method   string
}

func NewPaymentGateway(amount float64, currency string, method string) PaymentGateway {
	return &paymentGateway{
		amount:   amount,
		currency: currency,
		method:   method,
	}
}

func (pg *paymentGateway) ProcessPayment(ctx context.Context, userID string, amount float64) (bool, error) {
	return true, nil
}
