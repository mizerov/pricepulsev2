package elasticsearch

import (
	"testing"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	do.ProvideValue(i, &es8.Client{})

	repo, err := do.InvokeAs[domains.ProductSearchRepository](i)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}
