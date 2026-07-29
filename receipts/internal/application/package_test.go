package application

import (
	"testing"

	"github.com/mizerov/pricepulse/receipts/internal/application/commands"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	do.OverrideValue(i, &commands.UploadReceiptCommand{})
	do.OverrideValue(i, &commands.BindProductCommand{})

	app, err := do.Invoke[*Application](i)
	assert.NoError(t, err)
	assert.NotNil(t, app)
}
