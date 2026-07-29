package domains

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OutboxStatus string

const (
	OutboxStatusPending   OutboxStatus = "pending"
	OutboxStatusProcessed OutboxStatus = "processed"
)

type OutboxEvent struct {
	ID            uuid.UUID
	AggregateType string
	AggregateID   uuid.UUID
	Type          string
	Payload       []byte
	Status        OutboxStatus
	CreatedAt     time.Time
}

type OutboxEventRepository interface {
	FetchPendingEvents(ctx context.Context, batchSize int) ([]*OutboxEvent, error)
	MarkAsProcessed(ctx context.Context, id uuid.UUID) error
}
