package eventstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/EventStore/EventStore-Client-Go/v4/esdb"
)

// EventStoreDB is an implementation of EventStore using EventStoreDB.
type EventStoreDB struct {
	client *esdb.Client
}

// NewEventStoreDB creates a new EventStoreDB instance.
func NewEventStoreDB(connectionString string) (*EventStoreDB, error) {
	config, err := esdb.ParseConnectionString(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EventStoreDB connection string: %w", err)
	}

	client, err := esdb.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create EventStoreDB client: %w", err)
	}

	return &EventStoreDB{client: client}, nil
}

// Save persists a single event to the EventStoreDB.
func (e *EventStoreDB) Save(ctx context.Context, event Event) error {
	data, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	esEvent := esdb.EventData{
		ContentType: esdb.ContentTypeJson,
		EventType:   event.Type,
		Data:        data,
	}

	streamID := fmt.Sprintf("aggregate-%s", event.Aggregate)
	options := esdb.AppendToStreamOptions{}

	_, err = e.client.AppendToStream(ctx, streamID, options, esEvent)
	if err != nil {
		return fmt.Errorf("failed to save event to EventStoreDB: %w", err)
	}

	slog.Debug("saved event to EventStoreDB. stream id: %s", "streamID", streamID)

	return nil
}

// GetByAggregateID retrieves events for a specific aggregate ID from the EventStoreDB.
func (e *EventStoreDB) GetByAggregateID(ctx context.Context, aggregateID string) ([]Event, error) {
	streamID := fmt.Sprintf("aggregate-%s", aggregateID)
	options := esdb.ReadStreamOptions{}

	stream, err := e.client.ReadStream(ctx, streamID, options, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to read stream: %w", err)
	}
	defer stream.Close()

	var events []Event
	for {
		event, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("failed to receive event: %w", err)
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(event.OriginalEvent().Data, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event payload: %w", err)
		}

		events = append(events, Event{
			ID:        event.OriginalEvent().EventID.String(),
			Type:      event.OriginalEvent().EventType,
			Aggregate: aggregateID,
			Payload:   payload,
			Timestamp: event.OriginalEvent().CreatedDate,
		})
	}

	return events, nil
}
