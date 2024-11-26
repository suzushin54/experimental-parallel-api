package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	pb "github.com/suzushin54/experimental-parallel-api/gen/payment/v1"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/model"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/port"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/repository"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/gateway"
)

type ParallelPaymentService struct {
	paymentRepository repository.PaymentRepository
	paymentGateway    gateway.PaymentGateway
	idaasGateway      gateway.IDaaSInterface
	mailer            port.Mailer
	pb.UnimplementedPaymentServiceServer
}

func NewParallelPaymentService(
	pr repository.PaymentRepository,
	pg gateway.PaymentGateway,
	ig gateway.IDaaSInterface,
	m port.Mailer,
) *ParallelPaymentService {
	return &ParallelPaymentService{
		paymentRepository: pr,
		paymentGateway:    pg,
		idaasGateway:      ig,
		mailer:            m,
	}
}

func (s *ParallelPaymentService) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	paymentID, err := uuid.NewV7()
	if err != nil {
		return makeErrorResponse("Failed to generate UUID v7", err)
	}

	ptx, err := model.NewPaymentTransaction(paymentID.String(), req.PaymentData.Amount, req.PaymentData.Currency, req.PaymentData.Method)
	if err != nil {
		return makeErrorResponse("Transaction creation failed", err)
	}

	accountChan := make(chan string, 1)
	errChan := make(chan error, 2)

	go func() {
		id, err := s.idaasGateway.RegisterAccount(ctx, req.UserData.Email, req.UserData.Password)
		if err != nil {
			errChan <- err
			return
		}
		accountChan <- id
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

	close(accountChan)
	accountID := <-accountChan
	if err := ptx.BindCustomerToTransaction(accountID); err != nil {
		return makeErrorResponse("Transaction binding failed", err)
	}

	if err := s.paymentRepository.SaveTransaction(ctx, ptx); err != nil {
		return makeErrorResponse("Transaction saving failed", err)
	}

	if err := s.mailer.Send(ctx, req.UserData.Email, "Payment Confirmation", "Your payment has been processed successfully."); err != nil {
		return makeErrorResponse("Failed to send confirmation email", err)
	}

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
