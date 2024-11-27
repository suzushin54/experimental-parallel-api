package service

import (
	"context"

	"github.com/google/uuid"

	pb "github.com/suzushin54/experimental-parallel-api/gen/payment/v1"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/model"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/port"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/repository"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/gateway"
)

type SerialPaymentService struct {
	paymentRepository repository.PaymentRepository
	paymentGateway    gateway.PaymentGateway
	idaasGateway      gateway.IDaaSInterface
	mailer            port.Mailer
	pb.UnimplementedPaymentServiceServer
}

func NewSerialPaymentService(
	pr repository.PaymentRepository,
	pg gateway.PaymentGateway,
	ig gateway.IDaaSInterface,
	m port.Mailer,
) *SerialPaymentService {
	return &SerialPaymentService{
		paymentRepository: pr,
		paymentGateway:    pg,
		idaasGateway:      ig,
		mailer:            m,
	}
}

func (s *SerialPaymentService) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	paymentID, err := uuid.NewV7()
	if err != nil {
		return makeErrorResponse(ctx, "Failed to generate UUID v7", err)
	}

	ptx, err := model.NewPaymentTransaction(paymentID.String(), req.PaymentData.Amount, req.PaymentData.Currency, req.PaymentData.Method)
	if err != nil {
		return makeErrorResponse(ctx, "Transaction creation failed", err)
	}

	accountID, err := s.idaasGateway.RegisterAccount(ctx, req.UserData.Email, req.UserData.Password)
	if err != nil {
		return makeErrorResponse(ctx, "Account registration failed", err)
	}

	if err = s.paymentGateway.ProcessPayment(ctx, ptx); err != nil {
		return makeErrorResponse(ctx, "Payment processing failed", err)
	}

	if err = ptx.BindCustomerToTransaction(accountID); err != nil {
		return makeErrorResponse(ctx, "Transaction binding failed", err)
	}

	if err = s.paymentRepository.SaveTransaction(ctx, ptx); err != nil {
		return makeErrorResponse(ctx, "Transaction saving failed", err)
	}

	if err = s.mailer.Send(ctx, req.UserData.Email, "Payment Confirmation", "Your payment has been processed successfully"); err != nil {
		return makeErrorResponse(ctx, "Email sending failed", err)
	}

	return &pb.ProcessPaymentResponse{
		Success:      true,
		Message:      "Payment processed successfully",
		ErrorMessage: "",
	}, nil
}
