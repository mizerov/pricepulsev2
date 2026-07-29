package queries

import (
	"testing"

	"github.com/mizerov/pricepulse/catalog/internal/mocks"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	mockRepo := mocks.NewMockProductSearchRepository(t)
	do.ProvideValue[domains.ProductSearchRepository](i, mockRepo)

	searchQuery, err := do.Invoke[*SearchProductsQuery](i)
	assert.NoError(t, err)
	assert.NotNil(t, searchQuery)
}
