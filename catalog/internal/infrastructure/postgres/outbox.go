package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	tm "github.com/mizerov/pricepulse/catalog/pkg/transaction-manager"
	"github.com/samber/do/v2"
)

type OutboxRepository struct {
	db *pgxpool.Pool
}

func NewOutboxRepository(i do.Injector) (*OutboxRepository, error) {
	db, err := do.Invoke[*pgxpool.Pool](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke db: %w", err)
	}

	return &OutboxRepository{db: db}, nil
}

func (r *OutboxRepository) FetchPendingEvents(ctx context.Context, batchSize int) ([]*domains.OutboxEvent, error) {
	rows, err := tm.GetTxOrPool(ctx, r.db).Query(ctx,
		"SELECT id, aggregate_type, aggregate_id, type, payload, status, created_at FROM outbox_events WHERE status = 'pending' ORDER BY created_at ASC LIMIT $1 FOR UPDATE SKIP LOCKED", batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domains.OutboxEvent
	for rows.Next() {
		var e domains.OutboxEvent
		if err := rows.Scan(&e.ID, &e.AggregateType, &e.AggregateID, &e.Type, &e.Payload, &e.Status, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, &e)
	}
	return events, nil
}

func (r *OutboxRepository) MarkAsProcessed(ctx context.Context, id uuid.UUID) error {
	_, err := tm.GetTxOrPool(ctx, r.db).Exec(ctx, "UPDATE outbox_events SET status = 'processed' WHERE id = $1", id)
	return err
}
