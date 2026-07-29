package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mizerov/pricepulse/receipts/internal/application"
	"github.com/mizerov/pricepulse/receipts/internal/infrastructure/grpc"
	"github.com/mizerov/pricepulse/receipts/internal/infrastructure/postgres"
	"github.com/mizerov/pricepulse/receipts/internal/infrastructure/proverkacheka"
	httptransport "github.com/mizerov/pricepulse/receipts/internal/transport/http"
	"github.com/mizerov/router"
	"github.com/samber/do/v2"
)

type Config struct {
	DBHost             string `env:"DB_HOST" envDefault:"localhost"`
	DBPort             string `env:"DB_PORT" envDefault:"5434"`
	DBUser             string `env:"DB_USER" envDefault:"receipts"`
	DBPassword         string `env:"DB_PASSWORD" envDefault:"receipts"`
	DBName             string `env:"DB_NAME" envDefault:"receipts"`
	AppPort            string `env:"APP_PORT" envDefault:"8082"`
	JWTPublicKeyB64    string `env:"JWT_PUBLIC_KEY_B64"`
	ProverkaChekaToken string `env:"PROVERKACHEKA_TOKEN"`
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
		httptransport.Package,
		application.Package,
		postgres.Package,
		proverkacheka.Package,
		grpc.Package,
		func(i do.Injector) {
			do.ProvideValue(i, pool)
			do.ProvideValue(i, log)
			do.ProvideValue(i, proverkacheka.ProverkaCheckaConfig{
				HttpClient: &http.Client{},
				Token:      cfg.ProverkaChekaToken,
			})
			do.ProvideValue(i, httptransport.RouterConfig{PublicKeyB64: cfg.JWTPublicKeyB64})
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
