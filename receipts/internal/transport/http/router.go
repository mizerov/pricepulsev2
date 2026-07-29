package http

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mizerov/pricepulse/receipts/internal/application"

	"github.com/mizerov/router"
	"github.com/samber/do/v2"
)

type Handler struct {
	app *application.Application
}

type RouterConfig struct {
	PublicKeyB64 string
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

	routerConfig, err := do.Invoke[RouterConfig](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke router config: %w", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(routerConfig.PublicKeyB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 private key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	handler := &Handler{app: app}

	r := router.NewRouter(logger)

	router.AddRoute(r, http.MethodPost, "/api/receipts/v1/receipts/qrstring", handler.UploadReceiptByQRString, JWTMiddleware(publicKey))
	router.AddRoute(r, http.MethodPatch, "/api/receipts/v1/receipts/{id}/items/{item_id}", handler.BindProduct, JWTMiddleware(publicKey))

	return r, nil
}
