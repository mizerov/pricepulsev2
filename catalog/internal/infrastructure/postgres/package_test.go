package postgres

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	// Provide empty pool
	do.ProvideValue(i, &pgxpool.Pool{})

	productRepo, err := do.InvokeAs[domains.ProductRepository](i)
	assert.NoError(t, err)
	assert.NotNil(t, productRepo)

	outboxRepo, err := do.InvokeAs[domains.OutboxEventRepository](i)
	assert.NoError(t, err)
	assert.NotNil(t, outboxRepo)
}
