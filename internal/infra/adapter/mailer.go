package adapter

import (
	"context"
	"log/slog"
	"time"
)

// Mailer is a mock implementation of the Mailer interface.
type Mailer struct{}

// NewMailer creates a new instance of the Mailer.
func NewMailer() *Mailer {
	return &Mailer{}
}

// Send simulates sending an email and logs the operation.
func (m *Mailer) Send(ctx context.Context, to, subject, body string) error {
	slog.DebugContext(ctx, "Sending email to: %s", "to", to)
	slog.DebugContext(ctx, "Subject: %s", "subject", subject)
	slog.DebugContext(ctx, "Body: %s", "body", body)

	// Simulate network delay
	time.Sleep(1 * time.Second)

	slog.DebugContext(ctx, "Email successfully sent to: %s", "to", to)
	return nil
}
