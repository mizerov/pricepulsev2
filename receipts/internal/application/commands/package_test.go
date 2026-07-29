package commands

import (
	"testing"

	"github.com/mizerov/pricepulse/receipts/internal/mocks"
	"github.com/mizerov/pricepulse/receipts/pkg/domains"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	do.ProvideValue[domains.ReceiptRepository](i, mocks.NewMockReceiptRepository(t))
	do.ProvideValue[domains.CatalogClient](i, mocks.NewMockCatalogClient(t))
	do.ProvideValue[ReceiptProvider](i, mocks.NewMockReceiptProvider(t))

	uploadCmd, err := do.Invoke[*UploadReceiptCommand](i)
	assert.NoError(t, err)
	assert.NotNil(t, uploadCmd)

	bindCmd, err := do.Invoke[*BindProductCommand](i)
	assert.NoError(t, err)
	assert.NotNil(t, bindCmd)
}
