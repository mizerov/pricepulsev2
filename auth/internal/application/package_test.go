package application

import (
	"testing"

	"github.com/mizerov/pricepulse/auth/internal/application/commands"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	// Since application.Package includes commands.Package, 
	// they are already lazily provided. If we want to mock them, we must use OverrideValue.
	do.OverrideValue(i, &commands.CreateAccountCommand{})
	do.OverrideValue(i, &commands.IssueJWTCommand{})
	do.OverrideValue(i, &commands.RefreshJWTCommand{})

	app, err := do.Invoke[*Application](i)
	assert.NoError(t, err)
	assert.NotNil(t, app)
}
