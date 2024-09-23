package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Unlites/ml-analysis-provider/worker/internal/adapters/fulltextsearcher/elastic"
	natshandler "github.com/Unlites/ml-analysis-provider/worker/internal/adapters/handler/mq/nats"
	"github.com/Unlites/ml-analysis-provider/worker/internal/adapters/repository/postgres"
	"github.com/Unlites/ml-analysis-provider/worker/internal/application"
	"github.com/Unlites/ml-analysis-provider/worker/internal/config"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to create config", "detail", err)
		os.Exit(1)
	}

	ctx := context.Background()

	natsConn, err := nats.Connect(cfg.Nats.ConnString)
	if err != nil {
		slog.Error("failed to connect to nats", "detail", err)
		os.Exit(1)
	}
	defer natsConn.Drain()

	elasticConfig := elasticsearch.Config{
		Addresses: cfg.ElasticSearch.Addrs,
	}

	if cfg.ElasticSearch.CaPath != "" {
		cert, err := os.ReadFile(cfg.ElasticSearch.CaPath)
		if err != nil {
			slog.Error("failed to read elastic ca cert", "detail", err)
			os.Exit(1)
		}

		elasticConfig.CACert = cert
		elasticConfig.Username = cfg.ElasticSearch.Username
		elasticConfig.Password = cfg.ElasticSearch.Password
	}

	elasticClient, err := elasticsearch.NewTypedClient(elasticConfig)
	if err != nil {
		slog.Error("failed to connect to elasticsearch", "detail", err)
		os.Exit(1)
	}
	if success, err := elasticClient.Ping().IsSuccess(ctx); err != nil || !success {
		slog.Error("failed to ping elasticsearch", "detail", err)
		os.Exit(1)
	}

	pgPool, err := pgxpool.New(ctx, cfg.Postgres.ConnString)
	if err != nil {
		slog.Error("failed to create pg pool", "detail", err)
		os.Exit(1)
	}

	searcher := elastic.NewElasticFullTextSearcher(elasticClient)
	repo := postgres.NewPostgresRepository(pgPool)
	usecase := application.NewUsecase(searcher, repo)
	handler := natshandler.NewNatsMqHandler(natsConn, 5*time.Second, usecase)

	if err := handler.Start(); err != nil {
		slog.Error("failed to start handler", "detail", err)
		os.Exit(1)
	}

	notifyCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-notifyCtx.Done()

	if err := handler.Stop(); err != nil {
		slog.Error("failed to stop handler", "detail", err)
	}
}
