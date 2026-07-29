package application

import (
	"fmt"

	"github.com/mizerov/pricepulse/catalog/internal/application/commands"
	"github.com/mizerov/pricepulse/catalog/internal/application/queries"
	"github.com/samber/do/v2"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ProposeProduct *commands.ProposeProductCommand
	ApproveProduct *commands.ApproveProductCommand
}

type Queries struct {
	SearchProducts *queries.SearchProductsQuery
}

func NewApplication(i do.Injector) (*Application, error) {
	proposeProductCommand, err := do.Invoke[*commands.ProposeProductCommand](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke propose product command: %w", err)
	}
	approveProductCommand, err := do.Invoke[*commands.ApproveProductCommand](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke approve product command: %w", err)
	}
	searchProductsQuery, err := do.Invoke[*queries.SearchProductsQuery](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke search products query: %w", err)
	}

	return &Application{
		Commands: Commands{
			ProposeProduct: proposeProductCommand,
			ApproveProduct: approveProductCommand,
		},
		Queries: Queries{
			SearchProducts: searchProductsQuery,
		},
	}, nil
}
