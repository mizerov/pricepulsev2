package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/mizerov/pricepulse/catalog/internal/infrastructure/elasticsearch"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
)

type OutboxWorker struct {
	outboxRepo domains.OutboxEventRepository
	searchRepo *elasticsearch.ProductSearchRepository
	log        *slog.Logger
}

func NewOutboxWorker(i do.Injector) (*OutboxWorker, error) {
	outboxRepo, err := do.InvokeAs[domains.OutboxEventRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke outbox repo: %w", err)
	}
	searchRepo, err := do.Invoke[*elasticsearch.ProductSearchRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke search repo: %w", err)
	}
	log, err := do.Invoke[*slog.Logger](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke logger: %w", err)
	}

	return &OutboxWorker{
		outboxRepo: outboxRepo,
		searchRepo: searchRepo,
		log:        log,
	}, nil
}

func (w *OutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.log.Info("Outbox worker stopped")
			return
		case <-ticker.C:
			w.processEvents(ctx)
		}
	}
}

func (w *OutboxWorker) processEvents(ctx context.Context) {
	events, err := w.outboxRepo.FetchPendingEvents(ctx, 10)
	if err != nil {
		w.log.Error("Failed to fetch outbox events", "error", err)
		return
	}

	for _, event := range events {
		if event.AggregateType == "product" {
			var product domains.Product
			if err := json.Unmarshal(event.Payload, &product); err != nil {
				w.log.Error("Failed to unmarshal product from event", "error", err, "event_id", event.ID)
				continue
			}

			if err := w.searchRepo.IndexProduct(ctx, &product); err != nil {
				w.log.Error("Failed to index product", "error", err, "product_id", product.ID)
				continue
			}
		}

		if err := w.outboxRepo.MarkAsProcessed(ctx, event.ID); err != nil {
			w.log.Error("Failed to mark event as processed", "error", err, "event_id", event.ID)
		} else {
			w.log.Info("Processed outbox event", "event_id", event.ID)
		}
	}
}
