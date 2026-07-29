package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
)

type ProductSearchRepository struct {
	client *elasticsearch.Client
}

func NewProductSearchRepository(i do.Injector) (*ProductSearchRepository, error) {
	client, err := do.Invoke[*elasticsearch.Client](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke es client: %w", err)
	}
	return &ProductSearchRepository{client: client}, nil
}

func (r *ProductSearchRepository) CreateIndex(ctx context.Context) error {
	mapping := `{
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"name": { 
					"type": "text", 
					"analyzer": "standard" 
				},
				"description": { "type": "text" },
				"status": { "type": "keyword" }
			}
		}
	}`

	res, err := r.client.Indices.Create("products", r.client.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}
	return nil
}

func (r *ProductSearchRepository) IndexProduct(ctx context.Context, product *domains.Product) error {
	body, err := json.Marshal(map[string]interface{}{
		"id":          product.ID.String(),
		"name":        product.Name,
		"description": product.Description,
		"status":      string(product.Status),
	})
	if err != nil {
		return err
	}

	res, err := r.client.Index("products", bytes.NewReader(body), r.client.Index.WithDocumentID(product.ID.String()))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error indexing product: %s", res.String())
	}
	return nil
}

func (r *ProductSearchRepository) Search(ctx context.Context, query string) ([]*domains.Product, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"name": map[string]interface{}{
								"query":     query,
								"fuzziness": "AUTO",
							},
						},
					},
					{
						"term": map[string]interface{}{
							"status": string(domains.ProductStatusApproved),
						},
					},
				},
			},
		},
	}

	body, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("products"),
		r.client.Search.WithBody(bytes.NewReader(body)),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source struct {
					ID          uuid.UUID `json:"id"`
					Name        string    `json:"name"`
					Description string    `json:"description"`
					Status      string    `json:"status"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var products []*domains.Product
	for _, hit := range result.Hits.Hits {
		products = append(products, &domains.Product{
			ID:          hit.Source.ID,
			Name:        hit.Source.Name,
			Description: hit.Source.Description,
			Status:      domains.ProductStatus(hit.Source.Status),
		})
	}

	return products, nil
}
