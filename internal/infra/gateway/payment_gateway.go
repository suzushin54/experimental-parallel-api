package gateway

import (
	"context"
	"log/slog"
	"time"

	"github.com/suzushin54/experimental-parallel-api/internal/domain/model"
)

type PaymentGateway interface {
	ProcessPayment(ctx context.Context, ptx *model.PaymentTransaction) error
	ProcessPaymentWithDetails(ctx context.Context, paymentID string, amount float64, currency string, method string) error
}

type paymentGateway struct{}

func NewPaymentGateway() PaymentGateway {
	return &paymentGateway{}
}

func (pg *paymentGateway) ProcessPayment(ctx context.Context, ptx *model.PaymentTransaction) error {
	slog.DebugContext(ctx, "Processing payment transaction: %v", "ptx", ptx)

	// simulate a long-running transaction
	time.Sleep(800 * time.Millisecond)

	return nil
}

func (pg *paymentGateway) ProcessPaymentWithDetails(ctx context.Context, paymentID string, amount float64, currency string, method string) error {
	slog.DebugContext(ctx, "Processing payment with details: %v", "paymentID", paymentID, "amount", amount, "currency", currency, "method", method)

	// simulate a long-running transaction
	time.Sleep(800 * time.Millisecond)

	return nil
}
