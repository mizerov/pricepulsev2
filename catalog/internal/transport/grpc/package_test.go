package grpc

import (
	"testing"

	"github.com/mizerov/pricepulse/catalog/internal/mocks"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	mockProductRepo := mocks.NewMockProductRepository(t)

	i := do.New(Package)
	do.ProvideValue[domains.ProductRepository](i, mockProductRepo)

	server, err := do.Invoke[*Server](i)
	assert.NoError(t, err, "ошибка DI: не удалось инициализировать *Server")
	assert.NotNil(t, server, "ошибка DI: получен nil вместо *Server")
}
