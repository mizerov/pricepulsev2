package grpc

import (
	"context"

	catalogv1 "github.com/mizerov/pricepulse/api/catalog/v1"
	"google.golang.org/grpc"
)

type CatalogClient struct {
	client catalogv1.CatalogServiceClient
}

func NewCatalogClient(conn *grpc.ClientConn) *CatalogClient {
	return &CatalogClient{client: catalogv1.NewCatalogServiceClient(conn)}
}

func (c *CatalogClient) CheckProductExists(ctx context.Context, productID, accountID string) (bool, error) {
	resp, err := c.client.CheckProductExists(ctx, &catalogv1.CheckProductExistsRequest{
		ProductId: productID,
		AccountId: accountID,
	})
	if err != nil {
		return false, err
	}
	return resp.Exists, nil
}
