package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	pb "github.com/suzushin54/experimental-parallel-api/gen/payment/v1"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/port"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/eventstore"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/gateway"
)

// EventSourcedPaymentService is a service that processes payments using event sourcing.
type EventSourcedPaymentService struct {
	eventStore     eventstore.EventStore
	paymentGateway gateway.PaymentGateway
	idaasGateway   gateway.IDaaSInterface
	mailer         port.Mailer
	pb.UnimplementedPaymentServiceServer
}

// NewEventSourcedPaymentService creates a new EventSourcedPaymentService.
func NewEventSourcedPaymentService(
	es eventstore.EventStore,
	pg gateway.PaymentGateway,
	ig gateway.IDaaSInterface,
	m port.Mailer,
) *EventSourcedPaymentService {
	return &EventSourcedPaymentService{
		eventStore:     es,
		paymentGateway: pg,
		idaasGateway:   ig,
		mailer:         m,
	}
}

func (e *EventSourcedPaymentService) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	now := time.Now()
	paymentID, err := uuid.NewV7()
	if err != nil {
		return makeErrorResponse(ctx, "Failed to generate UUID v7", err)
	}

	initialEvent := eventstore.Event{
		ID:        uuid.NewString(),
		Type:      "PaymentInitiated",
		Aggregate: paymentID.String(),
		Payload: map[string]interface{}{
			"user": map[string]interface{}{
				"email":    req.UserData.Email,
				"password": req.UserData.Password, // TODO: need encryption
			},
			"payment": map[string]interface{}{
				"amount":   req.PaymentData.Amount,
				"currency": req.PaymentData.Currency,
				"method":   req.PaymentData.Method,
			},
		},
		Timestamp: now,
	}
	if err := e.eventStore.Save(ctx, initialEvent); err != nil {
		return makeErrorResponse(ctx, "Failed to save initial event", err)
	}

	// Channels for asynchronous processing
	accountChan := make(chan string, 1)
	errChan := make(chan error, 2)

	// IDaaS (Account Registration)
	go func() {
		accountID, err := e.idaasGateway.RegisterAccount(ctx, req.UserData.Email, req.UserData.Password)
		if err != nil {
			errChan <- fmt.Errorf("failed to register account: %w", err)
			return
		}

		accountEvent := eventstore.Event{
			ID:        uuid.NewString(),
			Type:      "AccountRegistered",
			Aggregate: paymentID.String(),
			Payload: map[string]interface{}{
				"accountID": accountID,
			},
			Timestamp: time.Now(),
		}
		if err := e.eventStore.Save(ctx, accountEvent); err != nil {
			errChan <- fmt.Errorf("failed to save account event: %w", err)
			return
		}

		accountChan <- accountID
		errChan <- nil
	}()

	// Payment Processing
	go func() {
		if err := e.paymentGateway.ProcessPaymentWithDetails(
			ctx,
			paymentID.String(),
			req.PaymentData.Amount,
			req.PaymentData.Currency,
			req.PaymentData.Method,
		); err != nil {
			errChan <- fmt.Errorf("failed to process payment: %w", err)
			return
		}

		// Save the payment processed event
		paymentEvent := eventstore.Event{
			ID:        uuid.NewString(),
			Type:      "PaymentProcessed",
			Aggregate: paymentID.String(),
			Payload: map[string]interface{}{
				"amount":   req.PaymentData.Amount,
				"currency": req.PaymentData.Currency,
				"method":   req.PaymentData.Method,
			},
			Timestamp: time.Now(),
		}
		if err := e.eventStore.Save(ctx, paymentEvent); err != nil {
			errChan <- fmt.Errorf("failed to save payment event: %w", err)
			return
		}

		errChan <- nil
	}()

	return &pb.ProcessPaymentResponse{
		Success:      true,
		Message:      "Payment initiated",
		ErrorMessage: "",
	}, nil
}
