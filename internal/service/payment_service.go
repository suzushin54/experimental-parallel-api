package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	pb "github.com/suzushin54/experimental-parallel-api/gen/payment/v1"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/model"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/repository"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/gateway"
)

type PaymentService struct {
	paymentRepository repository.PaymentRepository
	paymentGateway    gateway.PaymentGateway
	idaasGateway      gateway.IDaaSInterface
	pb.UnimplementedPaymentServiceServer
}

func NewPaymentService(
	pr repository.PaymentRepository,
	pg gateway.PaymentGateway,
	ig gateway.IDaaSInterface,
) *PaymentService {
	return &PaymentService{
		paymentRepository: pr,
		paymentGateway:    pg,
		idaasGateway:      ig,
	}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	log.Printf("ProcessPayment request received: %v", req)

	paymentID, err := uuid.NewV7()
	if err != nil {
		return makeErrorResponse("Failed to generate UUID v7", err)
	}

	ptx, err := model.NewPaymentTransaction(paymentID.String(), req.PaymentData.Amount, req.PaymentData.Currency, req.PaymentData.Method)
	if err != nil {
		return makeErrorResponse("Transaction creation failed", err)
	}

	errChan := make(chan error, 2)
	var accountID string

	go func() {
		id, err := s.idaasGateway.RegisterAccount(ctx, req.UserData.Email, req.UserData.Password)
		if err != nil {
			errChan <- err
			return
		}
		accountID = id
		log.Printf("Account registered: %s", id)
		errChan <- nil
	}()

	go func() {
		if err := s.paymentGateway.ProcessPayment(ctx, ptx); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return makeErrorResponse("Transaction processing failed", err)
		}
	}

	ptx.BindCustomerToTransaction(accountID)

	if err := s.paymentRepository.SaveTransaction(ctx, ptx); err != nil {
		return makeErrorResponse("Transaction saving failed", err)
	}

	log.Printf("Payment transaction created: %v", ptx)

	return &pb.ProcessPaymentResponse{
		Success:      true,
		Message:      "Payment processed successfully",
		ErrorMessage: "",
	}, nil
}

func makeErrorResponse(message string, err error) (*pb.ProcessPaymentResponse, error) {
	log.Printf("%s: %v", message, err)
	return &pb.ProcessPaymentResponse{
		Success:      false,
		Message:      "",
		ErrorMessage: err.Error(),
	}, nil
}
