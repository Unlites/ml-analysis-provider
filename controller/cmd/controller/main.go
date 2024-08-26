package main

import (
	"log/slog"
	"net/http"
	"os"

	natsbroker "github.com/Unlites/ml-analysis-provider/controller/internal/adapters/broker/nats"
	httphandler "github.com/Unlites/ml-analysis-provider/controller/internal/adapters/handler/http"
	"github.com/Unlites/ml-analysis-provider/controller/internal/application"
	"github.com/nats-io/nats.go"
)

func main() {
	natsConn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		slog.Error("failed to connect to nats", "detail", err)
		os.Exit(1)
	}

	broker := natsbroker.NewNatsBroker(natsConn)
	usecase := application.NewUsecase(broker)
	handler, err := httphandler.NewHTTPHandler(usecase)
	if err != nil {
		slog.Error("failed to create http handler", "detail", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("failed to start server", "detail", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
