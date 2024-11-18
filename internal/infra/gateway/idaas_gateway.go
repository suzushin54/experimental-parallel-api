package gateway

import (
	"context"
	"github.com/google/uuid"
	"log"
	"log/slog"
	"time"
)

// IDaaSInterface defines the interface for interacting with an IDaaS provider.
type IDaaSInterface interface {
	RegisterAccount(ctx context.Context, email, password string) (string, error)
}

// IDaaSGateway provides an interface to an external IDaaS provider.
type IDaaSGateway struct{}

// NewIDaaSGateway creates a new instance of an IDaaS service interface.
func NewIDaaSGateway() IDaaSInterface {
	return &IDaaSGateway{}
}

// RegisterAccount simulates registering a new account with an IDaaS provider.
func (i *IDaaSGateway) RegisterAccount(ctx context.Context, email, password string) (string, error) {
	slog.InfoContext(ctx, "Registering new account: %s", email)

	id, err := uuid.NewV7()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate UUID v7: %v", err)
		return "", err
	}

	log.Default().Printf("Account registered, ID: %s", id)

	// Simulate network delay
	time.Sleep(500 * time.Millisecond)

	return id.String(), nil
}
