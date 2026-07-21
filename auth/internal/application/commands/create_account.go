package commands

import (
	"context"
	"fmt"

	"github.com/mizerov/pricepulse/auth/pkg/domains"

	"github.com/samber/do/v2"
)

type CreateAccountCommand struct {
	accountRepository domains.AccountRepository
}

type CreateAccountData struct {
	Email    string
	Password string
}

func NewCreateAccountCommand(i do.Injector) (*CreateAccountCommand, error) {
	accountRepository, err := do.InvokeAs[domains.AccountRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke account repository: %w", err)
	}

	return &CreateAccountCommand{accountRepository: accountRepository}, nil
}

func (self *CreateAccountCommand) Execute(ctx context.Context, data CreateAccountData) error {
	account, err := domains.NewAccount(data.Email, data.Password)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	err = self.accountRepository.Create(ctx, account)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}
