package domain

import "time"

type Event struct {
	AggregateID string
	Type        string
	OccurredAt  time.Time
	Payload     map[string]any
}
