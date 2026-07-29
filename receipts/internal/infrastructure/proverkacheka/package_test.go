package proverkacheka

import (
	"net/http"
	"testing"

	"github.com/mizerov/pricepulse/receipts/internal/application/commands"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	do.ProvideValue(i, ProverkaCheckaConfig{
		HttpClient: &http.Client{},
		Token:      "test-token",
	})

	api, err := do.InvokeAs[commands.ReceiptProvider](i)
	assert.NoError(t, err)
	assert.NotNil(t, api)
}
