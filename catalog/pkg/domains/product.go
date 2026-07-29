package domains

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrForbidden       = errors.New("forbidden")
)

type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusApproved ProductStatus = "approved"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Status      ProductStatus
	OwnerID     uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Product) Approve() {
	p.Status = ProductStatusApproved
	p.UpdatedAt = time.Now()
}

type ProductRepository interface {
	CreateWithEvent(ctx context.Context, product *Product, event *OutboxEvent) error
	Update(ctx context.Context, product *Product) error
	UpdateWithEvent(ctx context.Context, product *Product, event *OutboxEvent) error
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
}

type ProductSearchRepository interface {
	Search(ctx context.Context, query string) ([]*Product, error)
}
