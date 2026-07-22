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

func NewRouter(i do.Injector) (*router.Router, error) {
	logger, err := do.Invoke[*slog.Logger](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke logger: %w", err)
	}

	app, err := do.Invoke[*application.Application](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke application: %w", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString("LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF5dnIvQTdQaEdVeFhQcnU3WFZXcwpuak1ZTS9JUEhZM0tPVEgxd0FmM2hpZmh1UnpkOFp2bkRib01ackM3UkttQzZrM1Voemo5SVQ4WU5oN2l1a2F1CkNEcUQxRjYwd3pCR08zdHU5bGVtNFJ5bWJPSEFQSUZ4OVBtb05kTXVHL1hBQ3NqVWsvYkdZVndwdXFldW1GSkwKdkpWc1pOMjJHSzhQdWxZOGJPNUcvKzg5RUNZTnoyNTZyWlpOSVRvSThBRjViSjRnTUVFUDRZaFNmQlRLSHl4TApkN0pGRklLdmdxYVdOUG1YckZCeXNzcGdJMkwzNHRBQnlOVXZCK2s0TjNyOUowOC9SVEVoUjU4UlJ5ODl1YnF1CmJWVHI2VDVDWUtSMEhGSUtLNTFGR00veEY2UW1WeFgzRnBEMVFVUjhwMmVaVjdsZE50alBFSlRUSmtaV3VNYnkKM3dJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg==\n")
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
	router.AddRoute(r, http.MethodPost, "/api/receipts/v1/products/offer", handler.OfferProduct, JWTMiddleware(publicKey))

	return r, nil
}
