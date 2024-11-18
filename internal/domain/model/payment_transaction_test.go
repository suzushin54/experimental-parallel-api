package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewPaymentTransaction(t *testing.T) {
	id, _ := uuid.NewV7()
	t.Run("Valid PaymentTransaction", func(t *testing.T) {
		pt, err := NewPaymentTransaction(id.String(), 100.0, "JPY", "credit")
		assert.NoError(t, err)
		assert.NotNil(t, pt)
		assert.Equal(t, "pending", pt.Status)
		assert.Equal(t, id.String(), pt.ID)
		assert.Equal(t, 100.0, pt.Amount)
		assert.Equal(t, "JPY", pt.Currency)
		assert.Equal(t, "credit", pt.PaymentMethod)
		assert.WithinDuration(t, time.Now(), pt.CreatedAt, time.Second)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		_, err := NewPaymentTransaction("invalid-uuid", 100.0, "USD", "credit")
		assert.Error(t, err)
	})

	t.Run("Negative Amount", func(t *testing.T) {
		_, err := NewPaymentTransaction(id.String(), -10.0, "USD", "credit")
		assert.Error(t, err)
	})

	t.Run("Invalid Currency", func(t *testing.T) {
		_, err := NewPaymentTransaction(id.String(), 100.0, "Ethereum", "credit")
		assert.Error(t, err)
	})

	t.Run("Invalid PaymentMethod", func(t *testing.T) {
		_, err := NewPaymentTransaction(id.String(), 100.0, "USD", "INVALID")
		assert.Error(t, err)
	})
}

func TestBindCustomerToTransaction(t *testing.T) {
	validID, _ := uuid.NewV7()
	validCustomerID, _ := uuid.NewV7()

	t.Run("Bind Customer to Transaction", func(t *testing.T) {
		pt, err := NewPaymentTransaction(validID.String(), 100.0, "USD", "credit")
		assert.NoError(t, err)

		err = pt.BindCustomerToTransaction(validCustomerID.String())
		assert.NoError(t, err)
		assert.Equal(t, validCustomerID.String(), *pt.CustomerID)
	})

	t.Run("Bind Customer Twice", func(t *testing.T) {
		pt, err := NewPaymentTransaction(validID.String(), 100.0, "USD", "credit")
		assert.NoError(t, err)

		err = pt.BindCustomerToTransaction(validCustomerID.String())
		assert.NoError(t, err)

		// Attempt to bind the customer ID again
		err = pt.BindCustomerToTransaction(validCustomerID.String())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "customer ID already bound to transaction")
	})
}
