package application

import (
	"github.com/mizerov/pricepulse/receipts/internal/application/commands"

	"github.com/samber/do/v2"
)

var Package = do.Package(
	commands.Package,
	do.Lazy(NewApplication),
)
