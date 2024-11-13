package service

import (
	"context"
	"log"

	"github.com/google/uuid"

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

func NewPaymentService(
	pr repository.PaymentRepository,
	pg gateway.PaymentGateway,
) *PaymentService {
	return &PaymentService{
		paymentRepository: pr,
		paymentGateway:    pg,
	}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("ProcessPayment request received: %v", req)

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	ptx, err := model.NewPaymentTransaction(id.String(), req.Amount, req.Currency, req.UserId, req.Method, "pending")
	if err != nil {
		return makeErrorResponse("Transaction creation failed", err)
	}

	if err := s.paymentGateway.ProcessPayment(ctx, ptx); err != nil {
		return makeErrorResponse("Payment processing failed", err)
	}

	if err := s.paymentRepository.SaveTransaction(ctx, ptx); err != nil {
		return makeErrorResponse("Transaction saving failed", err)
	}

	log.Printf("Payment transaction created: %v", ptx)

	return &pb.PaymentResponse{
		Success:       true,
		TransactionId: ptx.ID,
		Message:       "Payment processed successfully",
		ErrorMessage:  "",
	}, nil
}

func makeErrorResponse(message string, err error) (*pb.PaymentResponse, error) {
	log.Printf("%s: %v", message, err)
	return &pb.PaymentResponse{
		Success:       false,
		TransactionId: "",
		Message:       "",
		ErrorMessage:  err.Error(),
	}, nil
}
