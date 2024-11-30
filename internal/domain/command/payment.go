package command

import (
	"errors"
	"time"
)

// CommandType defines the type of command
type CommandType string

const (
	BindCustomer        CommandType = "BindCustomer"
	CompleteTransaction CommandType = "CompleteTransaction"
	FailTransaction     CommandType = "FailTransaction"
)

// PaymentCommand represents a command to modify a payment transaction
type PaymentCommand struct {
	Type      CommandType
	PaymentID string
	Params    map[string]interface{}
	Timestamp time.Time
}

// NewBindCustomerCommand creates a command to bind a customer to a transaction
func NewBindCustomerCommand(paymentID, customerID string) (*PaymentCommand, error) {
	if customerID == "" {
		return nil, errors.New("customerID cannot be empty")
	}
	return &PaymentCommand{
		Type:      BindCustomer,
		PaymentID: paymentID,
		Params: map[string]interface{}{
			"customerID": customerID,
		},
		Timestamp: time.Now(),
	}, nil
}

// NewCompleteTransactionCommand creates a command to complete a transaction
func NewCompleteTransactionCommand(paymentID string) (*PaymentCommand, error) {
	if paymentID == "" {
		return nil, errors.New("paymentID cannot be empty")
	}
	return &PaymentCommand{
		Type:      CompleteTransaction,
		PaymentID: paymentID,
		Timestamp: time.Now(),
	}, nil
}

// NewFailTransactionCommand creates a command to fail a transaction
func NewFailTransactionCommand(paymentID string) (*PaymentCommand, error) {
	if paymentID == "" {
		return nil, errors.New("paymentID cannot be empty")
	}
	return &PaymentCommand{
		Type:      FailTransaction,
		PaymentID: paymentID,
		Timestamp: time.Now(),
	}, nil
}
