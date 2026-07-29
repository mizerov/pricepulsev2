package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mizerov/pricepulse/receipts/pkg/domains"
	"github.com/samber/do/v2"
)

type BindProductCommand struct {
	receiptRepo   domains.ReceiptRepository
	catalogClient domains.CatalogClient
}

func NewBindProductCommand(i do.Injector) (*BindProductCommand, error) {
	receiptRepo, err := do.InvokeAs[domains.ReceiptRepository](i)
	if err != nil {
		return nil, err
	}
	catalogClient, err := do.InvokeAs[domains.CatalogClient](i)
	if err != nil {
		return nil, err
	}
	return &BindProductCommand{receiptRepo: receiptRepo, catalogClient: catalogClient}, nil
}

type BindProductData struct {
	AccountID uuid.UUID
	ReceiptID string
	ItemID    string
	ProductID string
}

func (c *BindProductCommand) Execute(ctx context.Context, data BindProductData) error {
	exists, err := c.catalogClient.CheckProductExists(ctx, data.ProductID, data.AccountID.String())
	if err != nil {
		return fmt.Errorf("failed to check product existence: %w", err)
	}
	if !exists {
		return errors.New("product does not exist or access denied")
	}

	receiptID, err := uuid.Parse(data.ReceiptID)
	if err != nil {
		return err
	}
	itemID, err := uuid.Parse(data.ItemID)
	if err != nil {
		return err
	}
	productID, err := uuid.Parse(data.ProductID)
	if err != nil {
		return err
	}

	// This assumes receiptRepo has a BindProductToItem method
	if err := c.receiptRepo.BindProductToItem(ctx, receiptID, itemID, productID); err != nil {
		return err
	}

	return nil
}
