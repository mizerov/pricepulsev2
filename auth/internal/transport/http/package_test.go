package http

import (
	"log/slog"
	"os"
	"testing"

	"github.com/mizerov/pricepulse/auth/internal/application"
	"github.com/mizerov/router"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	i := do.New(Package)

	do.ProvideValue(i, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	do.ProvideValue(i, &application.Application{})

	r, err := do.Invoke[*router.Router](i)
	assert.NoError(t, err)
	assert.NotNil(t, r)
}
