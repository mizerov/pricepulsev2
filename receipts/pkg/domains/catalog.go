package domains

import "context"

type CatalogClient interface {
	CheckProductExists(ctx context.Context, productID, accountID string) (bool, error)
}
