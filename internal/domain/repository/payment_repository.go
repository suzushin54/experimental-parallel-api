package repository

import (
	"context"

	"github.com/suzushin54/experimental-parallel-api/internal/domain/model"
)

type PaymentRepository interface {
	SaveTransaction(ctx context.Context, pt *model.PaymentTransaction) error
}
