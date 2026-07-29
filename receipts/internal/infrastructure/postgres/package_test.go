package postgres

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mizerov/pricepulse/receipts/pkg/domains"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	do.ProvideValue(i, &pgxpool.Pool{})

	repo, err := do.InvokeAs[domains.ReceiptRepository](i)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}
