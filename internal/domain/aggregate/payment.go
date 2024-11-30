package aggregate

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// PaymentTransactionAggregate represents the aggregate for payment transactions
type PaymentTransactionAggregate struct {
	ID            string
	Amount        float64
	Currency      string
	CustomerID    *string
	Status        string
	PaymentMethod string
	CreatedAt     time.Time
}

// NewPaymentTransactionAggregate initializes a new payment transaction aggregate
func NewPaymentTransactionAggregate(
	id string,
	amount float64,
	currency string,
	method string,
) (*PaymentTransactionAggregate, error) {
	pt := &PaymentTransactionAggregate{
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

// Apply initiates a state change in the aggregate based on the given action and parameters
func (pta *PaymentTransactionAggregate) Apply(action string, params map[string]interface{}) error {
	switch action {
	case "BindCustomer":
		return pta.bindCustomer(params)
	case "CompleteTransaction":
		return pta.completeTransaction()
	case "FailTransaction":
		return pta.failTransaction()
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

// bindCustomer binds a customer ID to the transaction
func (pta *PaymentTransactionAggregate) bindCustomer(params map[string]interface{}) error {
	cID, ok := params["customerID"].(string)
	if !ok || cID == "" {
		return errors.New("invalid customer ID")
	}

	if pta.CustomerID != nil {
		return fmt.Errorf("customer ID already bound to transaction: %s", *pta.CustomerID)
	}

	pta.CustomerID = &cID
	return nil
}

// completeTransaction sets the transaction status to "completed"
func (pta *PaymentTransactionAggregate) completeTransaction() error {
	if pta.Status != "pending" {
		return fmt.Errorf("transaction cannot be completed from status: %s", pta.Status)
	}
	pta.Status = "completed"
	return nil
}

// failTransaction sets the transaction status to "failed"
func (pta *PaymentTransactionAggregate) failTransaction() error {
	if pta.Status != "pending" {
		return fmt.Errorf("transaction cannot be failed from status: %s", pta.Status)
	}
	pta.Status = "failed"
	return nil
}
