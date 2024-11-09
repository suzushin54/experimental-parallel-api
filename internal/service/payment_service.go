package service

import (
	"context"
	pb "github.com/suzushin54/experimental-parallel-api/gen/payment"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/repository"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/gateway"
	"log"
)

type PaymentService struct {
	paymentRepository repository.PaymentRepository
	paymentGateway    gateway.PaymentGateway
	pb.UnimplementedPaymentServiceServer
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("ProcessPayment: %v", req)

	// TODO: validation

	return &pb.PaymentResponse{
		Success:       true,
		TransactionId: "123456789012",
		Message:       "Payment processed successfully",
		ErrorMessage:  "",
	}, nil
}
