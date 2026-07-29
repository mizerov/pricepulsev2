package domains

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ReceiptRepository interface {
	Exists(ctx context.Context, fn, fd, fp string) (bool, error)
	Save(ctx context.Context, receipt *Receipt) error
	BindProductToItem(ctx context.Context, receiptID, itemID, productID uuid.UUID) error
}

type Receipt struct {
	ID        uuid.UUID
	StoreName string
	StoreINN  string
	TotalSum  int64
	DateTime  time.Time
	FN        string
	FD        string
	FP        string
	Items     []ReceiptItem

	CreatedAt time.Time
	CreatedBy uuid.UUID
}

type ReceiptItem struct {
	ID        uuid.UUID
	ReceiptID uuid.UUID
	ProductID *uuid.UUID
	Name      string
	Price     int64
	Quantity  float64
	Amount    int64
}
