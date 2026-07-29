package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	tm "github.com/mizerov/pricepulse/catalog/pkg/transaction-manager"
	"github.com/samber/do/v2"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(i do.Injector) (*ProductRepository, error) {
	db, err := do.Invoke[*pgxpool.Pool](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke db: %w", err)
	}

	return &ProductRepository{db: db}, nil
}

func (r *ProductRepository) CreateWithEvent(ctx context.Context, product *domains.Product, event *domains.OutboxEvent) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, 
		"INSERT INTO products (id, name, description, status, owner_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		product.ID, product.Name, product.Description, product.Status, product.OwnerID, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}

	if event != nil {
		_, err = tx.Exec(ctx,
			"INSERT INTO outbox_events (id, aggregate_type, aggregate_id, type, payload, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			event.ID, event.AggregateType, event.AggregateID, event.Type, event.Payload, event.Status, event.CreatedAt)
		if err != nil {
			return fmt.Errorf("insert outbox event: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *ProductRepository) Update(ctx context.Context, product *domains.Product) error {
	_, err := tm.GetTxOrPool(ctx, r.db).Exec(ctx,
		"UPDATE products SET name = $1, description = $2, status = $3, updated_at = $4 WHERE id = $5",
		product.Name, product.Description, product.Status, product.UpdatedAt, product.ID)
	return err
}

func (r *ProductRepository) UpdateWithEvent(ctx context.Context, product *domains.Product, event *domains.OutboxEvent) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		"UPDATE products SET name = $1, description = $2, status = $3, updated_at = $4 WHERE id = $5",
		product.Name, product.Description, product.Status, product.UpdatedAt, product.ID)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}

	if event != nil {
		_, err = tx.Exec(ctx,
			"INSERT INTO outbox_events (id, aggregate_type, aggregate_id, type, payload, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			event.ID, event.AggregateType, event.AggregateID, event.Type, event.Payload, event.Status, event.CreatedAt)
		if err != nil {
			return fmt.Errorf("insert outbox event: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*domains.Product, error) {
	var p domains.Product
	err := tm.GetTxOrPool(ctx, r.db).QueryRow(ctx,
		"SELECT id, name, description, status, owner_id, created_at, updated_at FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domains.ErrProductNotFound
		}
		return nil, err
	}
	return &p, nil
}
