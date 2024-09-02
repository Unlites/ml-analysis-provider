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
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

func main() {
	ctx := context.Background()

	natsConn, err := nats.Connect("nats://nats:4222")
	if err != nil {
		slog.Error("failed to connect to nats", "detail", err)
		os.Exit(1)
	}
	defer natsConn.Drain()

	elasticClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	})
	if err != nil {
		slog.Error("failed to connect to elasticsearch", "detail", err)
		os.Exit(1)
	}

	pgPool, err := pgxpool.New(ctx, "postgres://postgres:postgres_pass@postgres:5432/ml_analysis?&pool_max_conns=10")
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
