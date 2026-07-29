package postgres

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mizerov/pricepulse/auth/pkg/domains"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	do.ProvideValue(i, &pgxpool.Pool{})

	accountRepo, err := do.InvokeAs[domains.AccountRepository](i)
	assert.NoError(t, err)
	assert.NotNil(t, accountRepo)

	sessionRepo, err := do.InvokeAs[domains.SessionRepository](i)
	assert.NoError(t, err)
	assert.NotNil(t, sessionRepo)
}
