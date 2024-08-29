package nats

import (
	"encoding/json"
	"log/slog"

	natsclient "github.com/nats-io/nats.go"
)

func sendResponse(m *natsclient.Msg, r Response) {
	if r.Error != nil {
		slog.Error("error while processing request", "detail", r.Error)
	}

	jsonResponse, err := json.Marshal(r)
	if err != nil {
		slog.Error("failed to marshal response", "detail", err)
		return
	}

	if err := m.Respond(jsonResponse); err != nil {
		slog.Error("failed to send response", "detail", err)
	}
}
