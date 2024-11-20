package service

import (
	"context"
	"testing"

	pb "github.com/suzushin54/experimental-parallel-api/gen/payment/v1"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/gateway"
	repositoryImpl "github.com/suzushin54/experimental-parallel-api/internal/infra/repository"
)

func setupServices() (*SerialPaymentService, *ParallelPaymentService) {
	repo := repositoryImpl.NewMemoryPaymentRepository()
	idaasGateway := gateway.NewIDaaSGateway()
	paymentGateway := gateway.NewPaymentGateway()

	serialService := NewSerialPaymentService(repo, paymentGateway, idaasGateway)
	parallelService := NewParallelPaymentService(repo, paymentGateway, idaasGateway)

	return serialService, parallelService
}

// BenchmarkSerialPaymentService benchmarks the serial payment service.
func BenchmarkSerialPaymentService(b *testing.B) {
	serialService, _ := setupServices()
	ctx := context.Background()
	req := &pb.ProcessPaymentRequest{
		UserData: &pb.UserData{
			Email:    "test@example.com",
			Password: "password123",
		},
		PaymentData: &pb.PaymentData{
			Amount:   1000,
			Currency: "USD",
			Method:   "card",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := serialService.ProcessPayment(ctx, req)
		if err != nil {
			b.Error(err)
		}
	}
}

// BenchmarkParallelPaymentService benchmarks the parallel payment service.
func BenchmarkParallelPaymentService(b *testing.B) {
	_, parallelService := setupServices()
	ctx := context.Background()
	req := &pb.ProcessPaymentRequest{
		UserData: &pb.UserData{
			Email:    "test@example.com",
			Password: "password123",
		},
		PaymentData: &pb.PaymentData{
			Amount:   1000,
			Currency: "USD",
			Method:   "card",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parallelService.ProcessPayment(ctx, req)
		if err != nil {
			b.Error(err)
		}
	}
}
