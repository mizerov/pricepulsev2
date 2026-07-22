package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mizerov/pricepulse/receipts/pkg/domains"
	"github.com/samber/do/v2"
)

type ReceiptRepository struct {
	db *pgxpool.Pool
}

func NewReceiptRepository(i do.Injector) (*ReceiptRepository, error) {
	db, err := do.Invoke[*pgxpool.Pool](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke sqlx.DB: %w", err)
	}
	return &ReceiptRepository{db: db}, nil
}

func (self *ReceiptRepository) Exists(ctx context.Context, fn, fd, fp string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM receipts WHERE fn=$1 AND fd=$2 AND fp=$3)`
	var exists bool

	err := self.db.QueryRow(ctx, query, fn, fd, fp).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check exists: %w", err)
	}

	return exists, nil
}

func (self *ReceiptRepository) Save(ctx context.Context, receipt *domains.Receipt) error {
	tx, err := self.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	receiptQuery := `
		INSERT INTO receipts (id, total_sum, date_time, store_name, store_inn, fn, fd, fp, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = tx.Exec(ctx, receiptQuery,
		receipt.ID,
		receipt.TotalSum,
		receipt.DateTime,
		receipt.StoreName,
		receipt.StoreINN,
		receipt.FN,
		receipt.FD,
		receipt.FP,
		receipt.CreatedBy,
		receipt.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert receipt: %w", err)
	}

	var values []string
	var args []any
	for _, item := range receipt.Items {
		values = append(values, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d)", len(args)+1, len(args)+2, len(args)+3, len(args)+4, len(args)+5, len(args)+6))
		args = append(args, item.ID, item.ReceiptID, item.Name, item.Amount, item.Price, item.Quantity)
	}

	query := fmt.Sprintf("INSERT INTO receipt_items (id, receipt_id, name, amount, price, quantity) VALUES %s",
		strings.Join(values, ","))
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert receipt items: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
