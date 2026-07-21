package domains

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
)

type TokenIssuer struct {
	sessionRepository SessionRepository
	tokenTTL          time.Duration
	refreshTokenTTL   time.Duration
	signedString      string
}

type TokenIssuerConfig struct {
	TokenTTL        time.Duration
	RefreshTokenTTL time.Duration
	SignedString    string
}

func NewTokenIssuer(i do.Injector) (*TokenIssuer, error) {
	sessionRepository, err := do.InvokeAs[SessionRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke session repository: %w", err)
	}
	config, err := do.InvokeAs[TokenIssuerConfig](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke token issuer config: %w", err)
	}

	return &TokenIssuer{
		sessionRepository: sessionRepository,
		tokenTTL:          config.TokenTTL,
		refreshTokenTTL:   config.RefreshTokenTTL,
		signedString:      config.SignedString,
	}, nil
}

type JWtClaims struct {
	jwt.RegisteredClaims
}

type TokenPair struct {
	SignedToken  string
	RefreshToken string
}

func (self *TokenIssuer) Issue(ctx context.Context, accountID uuid.UUID) (*TokenPair, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth",
			Subject:   accountID.String(),
			Audience:  jwt.ClaimStrings{"all"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(self.tokenTTL)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	})

	signedToken, err := token.SignedString([]byte(self.signedString))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	rawRefreshToken := uuid.New().String()

	session := &Session{
		ID:        uuid.New(),
		AccountID: accountID,
		TokenHash: HashToken(rawRefreshToken),
		IsBlocked: false,
		ExpiresAt: time.Now().Add(self.refreshTokenTTL),
		CreatedAt: time.Now(),
	}

	err = self.sessionRepository.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return &TokenPair{
		SignedToken:  signedToken,
		RefreshToken: rawRefreshToken,
	}, nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
