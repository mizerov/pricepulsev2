package http

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const AccountIDKey contextKey = "accountID"

func JWTMiddleware(publicKey *rsa.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error": "authorization header is required"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error": "invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			claims := &jwt.RegisteredClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return publicKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"error": "invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			accountID, err := uuid.Parse(claims.Subject)
			if err != nil {
				http.Error(w, `{"error": "invalid account id"}`, http.StatusUnauthorized)
			}

			ctx := context.WithValue(r.Context(), AccountIDKey, accountID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAccountIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	accountID, ok := ctx.Value(AccountIDKey).(uuid.UUID)
	return accountID, ok
}
