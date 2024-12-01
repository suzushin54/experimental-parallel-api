package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	pb "github.com/suzushin54/experimental-parallel-api/gen/payment/v1"
	"github.com/suzushin54/experimental-parallel-api/internal/domain/repository"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/eventstore"
)

// EventSourcedPaymentService is a service that processes payments using event sourcing.
type EventSourcedPaymentService struct {
	eventStore        eventstore.EventStore
	paymentRepository repository.PaymentRepository
	pb.UnimplementedPaymentServiceServer
}

// NewEventSourcedPaymentService creates a new EventSourcedPaymentService.
func NewEventSourcedPaymentService(
	es eventstore.EventStore,
	pr repository.PaymentRepository,
) *EventSourcedPaymentService {
	return &EventSourcedPaymentService{
		eventStore:        es,
		paymentRepository: pr,
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

	// TODO: implement each process triggered by the event

	return &pb.ProcessPaymentResponse{
		Success:      true,
		Message:      "Payment initiated",
		ErrorMessage: "",
	}, nil
}
