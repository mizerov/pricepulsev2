package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/mizerov/pricepulse/auth/internal/application"
	"github.com/mizerov/pricepulse/auth/internal/infrastructure/postgres"
	httptransport "github.com/mizerov/pricepulse/auth/internal/transport/http"
	"github.com/mizerov/pricepulse/auth/pkg/domains"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mizerov/router"
	"github.com/samber/do/v2"
)

type Config struct {
	DBHost          string        `env:"DB_HOST" envDefault:"localhost"`
	DBPort          string        `env:"DB_PORT" envDefault:"5432"`
	DBUser          string        `env:"DB_USER" envDefault:"auth"`
	DBPassword      string        `env:"DB_PASSWORD" envDefault:"auth"`
	DBName          string        `env:"DB_NAME" envDefault:"auth"`
	TokenTTL        time.Duration `env:"TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"336h"`
	SignedString    string        `env:"SIGNED_STRING" envDefault:"test"`
	AppPort         string        `env:"APP_PORT" envDefault:"8080"`
}

func main() {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Error("failed to parse env", "error", err)
		return
	}

	pool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		return
	}

	i := do.New(
		domains.Package,
		httptransport.Package,
		application.Package,
		postgres.Package,
		func(i do.Injector) {
			do.ProvideValue(i, pool)
			do.ProvideValue(i, log)
			do.ProvideValue(i, domains.TokenIssuerConfig{
				TokenTTL:        cfg.TokenTTL,
				RefreshTokenTTL: cfg.RefreshTokenTTL,
				SignedString:    cfg.SignedString,
			})
		},
	)

	router, err := do.Invoke[*router.Router](i)
	if err != nil {
		log.Error("failed to initialize router", "error", err)
		return
	}

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppPort), router)
	if err != nil {
		log.Error("failed to start server", "error", err)
		return
	}
}
