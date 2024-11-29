package eventstore

import "context"

// EventStore is an interface for persisting and retrieving domain events.
type EventStore interface {
	Save(ctx context.Context, event Event) error
	GetByAggregateID(ctx context.Context, aggregateID string) ([]Event, error)
}
