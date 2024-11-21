package gateway

import (
	"context"
	"log/slog"
	"time"

	"github.com/suzushin54/experimental-parallel-api/internal/domain/model"
)

type PaymentGateway interface {
	ProcessPayment(ctx context.Context, ptx *model.PaymentTransaction) error
}

type paymentGateway struct{}

func NewPaymentGateway() PaymentGateway {
	return &paymentGateway{}
}

func (pg *paymentGateway) ProcessPayment(ctx context.Context, ptx *model.PaymentTransaction) error {
	slog.DebugContext(ctx, "Processing payment transaction: %v", ptx)

	// simulate a long-running transaction
	time.Sleep(800 * time.Millisecond)

	return nil
}
