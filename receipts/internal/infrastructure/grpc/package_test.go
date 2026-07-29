package grpc

import (
	"testing"

	"github.com/mizerov/pricepulse/receipts/pkg/domains"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	// Here it invokes standard gRPC dial under the hood. Since it's Lazy and the func dials localhost:9090.
	// It doesn't actually connect fully because insecure dial without block is non-blocking.
	// But it might fail if we are strict. Let's just invoke it.
	client, err := do.InvokeAs[domains.CatalogClient](i)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
