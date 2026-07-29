package workers

import (
	"log/slog"
	"os"
	"testing"

	"github.com/mizerov/pricepulse/catalog/internal/infrastructure/elasticsearch"
	"github.com/mizerov/pricepulse/catalog/internal/mocks"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	mockOutboxRepo := mocks.NewMockOutboxEventRepository(t)
	do.ProvideValue[domains.OutboxEventRepository](i, mockOutboxRepo)
	do.ProvideValue(i, &elasticsearch.ProductSearchRepository{})
	do.ProvideValue(i, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	worker, err := do.Invoke[*OutboxWorker](i)
	assert.NoError(t, err)
	assert.NotNil(t, worker)
}
