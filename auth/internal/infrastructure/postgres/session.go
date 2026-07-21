package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/mizerov/pricepulse/auth/pkg/domains"
	tm "github.com/mizerov/pricepulse/auth/pkg/transaction-manager"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(i do.Injector) (*SessionRepository, error) {
	db, err := do.Invoke[*pgxpool.Pool](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke db: %w", err)
	}

	return &SessionRepository{db: db}, nil
}

func (self *SessionRepository) Create(ctx context.Context, session *domains.Session) error {
	_, err := tm.GetTxOrPool(ctx, self.db).
		Exec(ctx, "INSERT INTO sessions (id, account_id, token_hash, is_blocked, expires_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
			session.ID, session.AccountID, session.TokenHash, session.IsBlocked, session.ExpiresAt, session.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (self *SessionRepository) FindByHash(ctx context.Context, tokenHash string) (*domains.Session, error) {
	var session domains.Session

	err := tm.GetTxOrPool(ctx, self.db).
		QueryRow(ctx, "SELECT id, account_id, token_hash, is_blocked, expires_at, created_at FROM sessions WHERE token_hash = $1", tokenHash).
		Scan(&session.ID, &session.AccountID, &session.TokenHash, &session.IsBlocked, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domains.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &session, nil
}

func (self *SessionRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := tm.GetTxOrPool(ctx, self.db).Exec(ctx, "DELETE FROM sessions WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
