package grpc

import (
	"github.com/mizerov/pricepulse/receipts/pkg/domains"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Package = do.Package(
	do.Lazy(func(i do.Injector) (domains.CatalogClient, error) {
		conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials())) // Note: should take address from config in a real app
		if err != nil {
			return nil, err
		}
		return NewCatalogClient(conn), nil
	}),
)
