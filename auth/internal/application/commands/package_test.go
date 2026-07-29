package commands

import (
	"testing"

	"github.com/mizerov/pricepulse/auth/internal/mocks"
	"github.com/mizerov/pricepulse/auth/pkg/domains"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	do.ProvideValue[domains.AccountRepository](i, mocks.NewMockAccountRepository(t))
	do.ProvideValue[domains.SessionRepository](i, mocks.NewMockSessionRepository(t))
	do.ProvideValue(i, &domains.TokenIssuer{})

	createCmd, err := do.Invoke[*CreateAccountCommand](i)
	assert.NoError(t, err)
	assert.NotNil(t, createCmd)

	issueCmd, err := do.Invoke[*IssueJWTCommand](i)
	assert.NoError(t, err)
	assert.NotNil(t, issueCmd)

	refreshCmd, err := do.Invoke[*RefreshJWTCommand](i)
	assert.NoError(t, err)
	assert.NotNil(t, refreshCmd)
}
