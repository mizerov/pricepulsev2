package queries

import (
	"context"

	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
)

type SearchProductsQuery struct {
	searchRepo domains.ProductSearchRepository
}

func NewSearchProductsQuery(i do.Injector) (*SearchProductsQuery, error) {
	searchRepo, err := do.InvokeAs[domains.ProductSearchRepository](i)
	if err != nil {
		return nil, err
	}
	return &SearchProductsQuery{searchRepo: searchRepo}, nil
}

func (q *SearchProductsQuery) Execute(ctx context.Context, query string) ([]*domains.Product, error) {
	return q.searchRepo.Search(ctx, query)
}
