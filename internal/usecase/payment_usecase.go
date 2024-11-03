package usecase

import (
	"context"

	"github.com/suzushin54/experimental-parallel-api/internal/infra/gateway"
)

type PaymentUseCase struct {
	paymentGateway gateway.PaymentGateway
}

// NewPaymentUseCase はPaymentUseCaseの新しいインスタンスを生成します。
func NewPaymentUseCase(pg gateway.PaymentGateway) *PaymentUseCase {
	return &PaymentUseCase{
		paymentGateway: pg,
	}
}

// Execute は決済処理のユースケースを実行します。
func (uc *PaymentUseCase) Execute(ctx context.Context, userID string, amount float64) error {
	uc.paymentGateway.ProcessPayment(ctx, "", 0.0)
	return nil
}
