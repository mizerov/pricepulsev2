package commands

import (
	"testing"

	"github.com/mizerov/pricepulse/catalog/internal/mocks"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	mockRepo := mocks.NewMockProductRepository(t)
	do.ProvideValue[domains.ProductRepository](i, mockRepo)

	proposeCmd, err := do.Invoke[*ProposeProductCommand](i)
	assert.NoError(t, err)
	assert.NotNil(t, proposeCmd)

	approveCmd, err := do.Invoke[*ApproveProductCommand](i)
	assert.NoError(t, err)
	assert.NotNil(t, approveCmd)
}
