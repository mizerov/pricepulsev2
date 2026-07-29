package application

import (
	"github.com/mizerov/pricepulse/catalog/internal/application/commands"
	"github.com/mizerov/pricepulse/catalog/internal/application/queries"

	"github.com/samber/do/v2"
)

var Package = do.Package(
	commands.Package,
	queries.Package,
	do.Lazy(NewApplication),
)
