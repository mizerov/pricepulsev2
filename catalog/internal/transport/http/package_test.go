package http

import (
	"testing"

	"github.com/mizerov/pricepulse/catalog/internal/application/commands"
	"github.com/mizerov/pricepulse/catalog/internal/application/queries"
	"github.com/mizerov/router"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	do.ProvideValue(i, &queries.SearchProductsQuery{})
	do.ProvideValue(i, &commands.ProposeProductCommand{})
	do.ProvideValue(i, &commands.ApproveProductCommand{})

	handlers, err := do.Invoke[*router.Router](i)
	assert.NoError(t, err)
	assert.NotNil(t, handlers)
}
