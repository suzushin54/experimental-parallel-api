package repository

import (
	"context"
	"time"

	"github.com/suzushin54/experimental-parallel-api/internal/domain/model"
)

// MemoryPaymentRepository implements PaymentRepository interface as mock.
type MemoryPaymentRepository struct {
	Transactions map[string]*model.PaymentTransaction
}

// NewMemoryPaymentRepository creates a new MemoryPaymentRepository.
func NewMemoryPaymentRepository() *MemoryPaymentRepository {
	return &MemoryPaymentRepository{
		Transactions: make(map[string]*model.PaymentTransaction),
	}
}

// SaveTransaction saves a payment transaction to memory.
func (repo *MemoryPaymentRepository) SaveTransaction(ctx context.Context, pt *model.PaymentTransaction) error {
	// simulate a long-running transaction
	time.Sleep(500 * time.Millisecond)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		repo.Transactions[pt.ID] = pt
		return nil
	}
}
