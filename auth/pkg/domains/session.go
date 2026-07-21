package domains

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrSessionNotFound = errors.New("session not found")
var ErrSessionIsBlocked = errors.New("session is blocked")

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	FindByHash(ctx context.Context, hash string) (*Session, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type Session struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	TokenHash string
	IsBlocked bool
	ExpiresAt time.Time
	CreatedAt time.Time
}
