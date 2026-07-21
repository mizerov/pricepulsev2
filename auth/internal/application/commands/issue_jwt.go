package commands

import (
	"context"
	"fmt"

	"github.com/mizerov/pricepulse/auth/pkg/domains"

	"github.com/samber/do/v2"
)

type IssueJWTCommand struct {
	accountRepository domains.AccountRepository
	tokenIssuer       *domains.TokenIssuer
}

type IssueJWTData struct {
	Email    string
	Password string
}

type IssueJWTResponse struct {
	SignedToken  string
	RefreshToken string
}

func NewIssueJWTCommand(i do.Injector) (*IssueJWTCommand, error) {
	accountRepository, err := do.InvokeAs[domains.AccountRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke account repository: %w", err)
	}
	tokenIssuer, err := do.Invoke[*domains.TokenIssuer](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke token issuer: %w", err)
	}

	return &IssueJWTCommand{accountRepository: accountRepository, tokenIssuer: tokenIssuer}, nil
}

func (self *IssueJWTCommand) Execute(ctx context.Context, data IssueJWTData) (*IssueJWTResponse, error) {
	account, err := self.accountRepository.FindByEmail(ctx, data.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find account by email: %w", err)
	}

	err = account.VerifyPassword(data.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to verify password: %w", err)
	}

	tokenPair, err := self.tokenIssuer.Issue(ctx, account.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to issue token: %w", err)
	}

	return &IssueJWTResponse{SignedToken: tokenPair.SignedToken, RefreshToken: tokenPair.RefreshToken}, nil
}
