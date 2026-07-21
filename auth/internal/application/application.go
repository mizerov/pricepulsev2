package application

import (
	"fmt"

	"github.com/mizerov/pricepulse/auth/internal/application/commands"

	"github.com/samber/do/v2"
)

type Application struct {
	Commands Commands
}

type Commands struct {
	CreateAccount *commands.CreateAccountCommand
	IssueJWT      *commands.IssueJWTCommand
	RefreshJWT    *commands.RefreshJWTCommand
}

func NewApplication(i do.Injector) (*Application, error) {
	createAccountCommand, err := do.Invoke[*commands.CreateAccountCommand](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke create account command: %w", err)
	}
	issueJWTCommand, err := do.Invoke[*commands.IssueJWTCommand](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke issue jwt command: %w", err)
	}
	refreshJWTCommand, err := do.Invoke[*commands.RefreshJWTCommand](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke refresh jwt command: %w", err)
	}

	return &Application{
		Commands: Commands{
			CreateAccount: createAccountCommand,
			IssueJWT:      issueJWTCommand,
			RefreshJWT:    refreshJWTCommand,
		},
	}, nil
}
