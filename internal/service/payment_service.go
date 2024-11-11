package service

import (
	"context"
	"log"

	pb "github.com/suzushin54/experimental-parallel-api/gen/payment"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/model"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/repository"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/gateway"
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

	tx, err := model.NewPaymentTransaction("generate-id", req.Amount, req.Currency, req.UserId, req.Method, "pending")
	if err != nil {
		log.Printf("Error creating payment transaction: %v", err)
		return &pb.PaymentResponse{
			Success:       false,
			TransactionId: "",
			Message:       "",
			ErrorMessage:  err.Error(),
		}, nil
	}

	if err := s.paymentRepository.SaveTransaction(ctx, tx); err != nil {
		log.Printf("Error saving payment transaction: %v", err)
		return &pb.PaymentResponse{
			Success:       false,
			TransactionId: "",
			Message:       "",
			ErrorMessage:  err.Error(),
		}, nil
	}

	log.Printf("Payment transaction created: %v", tx)

	return &pb.PaymentResponse{
		Success:       true,
		TransactionId: tx.ID,
		Message:       "Payment processed successfully",
		ErrorMessage:  "",
	}, nil
}
