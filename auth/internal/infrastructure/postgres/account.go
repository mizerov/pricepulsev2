package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mizerov/pricepulse/auth/pkg/domains"
	tm "github.com/mizerov/pricepulse/auth/pkg/transaction-manager"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(i do.Injector) (*AccountRepository, error) {
	db, err := do.Invoke[*pgxpool.Pool](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke db: %w", err)
	}

	return &AccountRepository{db: db}, nil
}

func (self *AccountRepository) FindByEmail(ctx context.Context, email string) (*domains.Account, error) {
	var account domains.Account

	err := tm.GetTxOrPool(ctx, self.db).
		QueryRow(ctx, "SELECT id, email, password_hash, role FROM accounts WHERE email = $1", email).
		Scan(&account.ID, &account.Email, &account.PasswordHash, &account.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domains.ErrAccountNotFound
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &account, nil
}

func (self *AccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*domains.Account, error) {
	var account domains.Account

	err := tm.GetTxOrPool(ctx, self.db).
		QueryRow(ctx, "SELECT id, email, password_hash, role FROM accounts WHERE id = $1", id).
		Scan(&account.ID, &account.Email, &account.PasswordHash, &account.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domains.ErrAccountNotFound
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &account, nil
}

func (self *AccountRepository) Create(ctx context.Context, account *domains.Account) error {
	_, err := tm.GetTxOrPool(ctx, self.db).
		Exec(ctx, "INSERT INTO accounts (id, email, password_hash, role) VALUES ($1, $2, $3, $4)",
			account.ID, account.Email, account.PasswordHash, account.Role)
	if err != nil {
		if tErr, ok := errors.AsType[*pgconn.PgError](err); ok && tErr.Code == "23505" {
			return domains.ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
