package commands

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
)

type ProposeProductCommand struct {
	productRepo domains.ProductRepository
}

func NewProposeProductCommand(i do.Injector) (*ProposeProductCommand, error) {
	productRepo, err := do.InvokeAs[domains.ProductRepository](i)
	if err != nil {
		return nil, err
	}
	return &ProposeProductCommand{productRepo: productRepo}, nil
}

type ProposeProductData struct {
	AccountID   uuid.UUID
	Name        string
	Description string
}

func (c *ProposeProductCommand) Execute(ctx context.Context, data ProposeProductData) (*domains.Product, error) {
	product := &domains.Product{
		ID:          uuid.New(),
		Name:        data.Name,
		Description: data.Description,
		Status:      domains.ProductStatusDraft,
		OwnerID:     data.AccountID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	payload, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}

	event := &domains.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "product",
		AggregateID:   product.ID,
		Type:          "ProductProposed",
		Payload:       payload,
		Status:        domains.OutboxStatusPending,
		CreatedAt:     time.Now(),
	}

	if err := c.productRepo.CreateWithEvent(ctx, product, event); err != nil {
		return nil, err
	}

	return product, nil
}
