package commands

import (
	"context"
	"fmt"

	"github.com/mizerov/pricepulse/auth/pkg/domains"

	"github.com/samber/do/v2"
)

type RefreshJWTCommand struct {
	sessionRepository domains.SessionRepository
	tokenIssuer       *domains.TokenIssuer
}

type RefreshJWTData struct {
	RefreshToken string
}

type RefreshJWTResponse struct {
	SignedToken  string
	RefreshToken string
}

func NewRefreshJWTCommand(i do.Injector) (*RefreshJWTCommand, error) {
	sessionRepository, err := do.InvokeAs[domains.SessionRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke session repository: %w", err)
	}
	tokenIssuer, err := do.Invoke[*domains.TokenIssuer](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke token issuer: %w", err)
	}

	return &RefreshJWTCommand{sessionRepository: sessionRepository, tokenIssuer: tokenIssuer}, nil
}

func (self *RefreshJWTCommand) Execute(ctx context.Context, data RefreshJWTData) (*RefreshJWTResponse, error) {
	session, err := self.sessionRepository.FindByHash(ctx, domains.HashToken(data.RefreshToken))
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	if session.IsBlocked {
		return nil, domains.ErrSessionIsBlocked
	}

	err = self.sessionRepository.DeleteByID(ctx, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete session: %w", err)
	}

	tokenPair, err := self.tokenIssuer.Issue(ctx, session.AccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to issue token: %w", err)
	}

	return &RefreshJWTResponse{
		SignedToken:  tokenPair.SignedToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
