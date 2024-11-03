package model

import "time"

type PaymentTransaction struct {
	ID            string
	Amount        float64
	Currency      string
	CustomerID    string
	Status        string
	PaymentMethod string
	CreatedAt     time.Time
}
