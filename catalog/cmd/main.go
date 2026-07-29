package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/caarlos0/env/v11"
	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	catalogv1 "github.com/mizerov/pricepulse/api/catalog/v1"
	"github.com/mizerov/pricepulse/catalog/internal/application"
	"github.com/mizerov/pricepulse/catalog/internal/application/workers"
	"github.com/mizerov/pricepulse/catalog/internal/infrastructure/elasticsearch"
	"github.com/mizerov/pricepulse/catalog/internal/infrastructure/postgres"
	grpcserver "github.com/mizerov/pricepulse/catalog/internal/transport/grpc"
	httptransport "github.com/mizerov/pricepulse/catalog/internal/transport/http"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/mizerov/router"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"
)

type Config struct {
	DBHost          string `env:"DB_HOST" envDefault:"localhost"`
	DBPort          string `env:"DB_PORT" envDefault:"5433"`
	DBUser          string `env:"DB_USER" envDefault:"catalog"`
	DBPassword      string `env:"DB_PASSWORD" envDefault:"catalog"`
	DBName          string `env:"DB_NAME" envDefault:"catalog"`
	ESUrl           string `env:"ES_URL" envDefault:"http://localhost:9200"`
	AppPort         string `env:"APP_PORT" envDefault:"8083"`
	JWTPublicKeyB64 string `env:"JWT_PUBLIC_KEY_B64"`
	GrpcPort        string `env:"GRPC_PORT" envDefault:"9090"`
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

	esClient, err := es8.NewClient(es8.Config{
		Addresses: []string{cfg.ESUrl},
	})
	if err != nil {
		log.Error("failed to create es client", "error", err)
		return
	}

	i := do.New(
		domains.Package,
		application.Package,
		workers.Package,
		postgres.Package,
		elasticsearch.Package,
		httptransport.Package,
		grpcserver.Package,
		func(i do.Injector) {
			do.ProvideValue(i, log)
			do.ProvideValue(i, pool)
			do.ProvideValue(i, esClient)
		},
	)

	// Start Outbox Worker
	worker, err := do.Invoke[*workers.OutboxWorker](i)
	if err != nil {
		log.Error("failed to invoke worker", "error", err)
		return
	}
	go worker.Start(ctx)

	// Start GRPC Server
	grpcSrv, err := do.Invoke[*grpcserver.Server](i)
	if err != nil {
		log.Error("failed to invoke grpc server", "error", err)
		return
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GrpcPort))
	if err != nil {
		log.Error("failed to listen grpc", "error", err)
		return
	}
	s := grpc.NewServer()
	catalogv1.RegisterCatalogServiceServer(s, grpcSrv)
	go func() {
		log.Info("Starting gRPC server", "port", cfg.GrpcPort)
		if err := s.Serve(lis); err != nil {
			log.Error("failed to serve grpc", "error", err)
		}
	}()

	// Start HTTP Server
	r, err := do.Invoke[*router.Router](i)
	if err != nil {
		log.Error("failed to invoke http handlers", "error", err)
		return
	}

	log.Info("Catalog HTTP server started", "port", cfg.AppPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppPort), r)
	if err != nil {
		log.Error("failed to start http server", "error", err)
		return
	}
}
