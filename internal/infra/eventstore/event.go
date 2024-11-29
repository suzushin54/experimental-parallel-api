package eventstore

import "time"

// Event represents a single domain event in the event store
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Aggregate string                 `json:"aggregate"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
}
