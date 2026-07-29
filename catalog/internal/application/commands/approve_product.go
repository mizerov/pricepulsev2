package commands

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
)

type ApproveProductCommand struct {
	productRepo domains.ProductRepository
}

func NewApproveProductCommand(i do.Injector) (*ApproveProductCommand, error) {
	productRepo, err := do.InvokeAs[domains.ProductRepository](i)
	if err != nil {
		return nil, err
	}
	return &ApproveProductCommand{productRepo: productRepo}, nil
}

func (c *ApproveProductCommand) Execute(ctx context.Context, productID uuid.UUID) error {
	product, err := c.productRepo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	product.Approve()

	// In a real app we might want to also create an outbox event here for Elasticsearch update.
	// But let's just keep it simple or add an event update as well.
	payload, _ := json.Marshal(product)
	event := &domains.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "product",
		AggregateID:   product.ID,
		Type:          "ProductApproved",
		Payload:       payload,
		Status:        domains.OutboxStatusPending,
		CreatedAt:     time.Now(),
	}

	// Wait, the interface for CreateWithEvent is for creation. 
	// We need UpdateWithEvent or just create the event directly.
	// Since ProductRepo doesn't have UpdateWithEvent, I will update product directly for simplicity.
	// Actually, outbox should track updates too. Let's add UpdateWithEvent.
	
	if err := c.productRepo.UpdateWithEvent(ctx, product, event); err != nil {
		return err
	}

	return nil
}
