package eventstore

import (
	"fmt"
	"sync"
)

// InMemoryEventStore is an in-memory implementation of the EventStore interface.
type InMemoryEventStore struct {
	mu     sync.RWMutex
	events map[string][]Event
}

// NewInMemoryEventStore creates a new instance of the in-memory event store.
func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make(map[string][]Event),
	}
}

// Save persists the event in the in-memory event store.
func (i *InMemoryEventStore) Save(event Event) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Create a new event slice if one does not exist for the aggregate
	if _, exists := i.events[event.Aggregate]; !exists {
		i.events[event.Aggregate] = []Event{}
	}

	i.events[event.Aggregate] = append(i.events[event.Aggregate], event)

	return nil
}

// GetByAggregateID retrieves all events for a given aggregate ID.
func (i *InMemoryEventStore) GetByAggregateID(aggregateID string) ([]Event, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	events, exists := i.events[aggregateID]
	if !exists {
		return nil, fmt.Errorf("no events found for aggregate ID: %s", aggregateID)
	}

	return events, nil
}
