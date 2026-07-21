package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mizerov/pricepulse/auth/internal/application"

	"github.com/mizerov/router"
	"github.com/samber/do/v2"
)

type Handler struct {
	app *application.Application
}

func NewRouter(i do.Injector) (*router.Router, error) {
	logger, err := do.Invoke[*slog.Logger](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke logger: %w", err)
	}

	app, err := do.Invoke[*application.Application](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke application: %w", err)
	}

	handler := &Handler{app: app}

	r := router.NewRouter(logger)

	router.AddRoute(r, http.MethodPost, "/api/auth/v1/jwt/issue", handler.IssueJWT)
	router.AddRoute(r, http.MethodPost, "/api/auth/v1/jwt/refresh", handler.RefreshJWT)
	router.AddRoute(r, http.MethodPost, "/api/auth/v1/accounts", handler.CreateAccount)

	return r, nil
}
