package application

import (
	"testing"

	"github.com/mizerov/pricepulse/catalog/internal/application/commands"
	"github.com/mizerov/pricepulse/catalog/internal/application/queries"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	do.OverrideValue(i, &commands.ProposeProductCommand{})
	do.OverrideValue(i, &commands.ApproveProductCommand{})
	do.OverrideValue(i, &queries.SearchProductsQuery{})

	app, err := do.Invoke[*Application](i)
	assert.NoError(t, err)
	assert.NotNil(t, app)
}
