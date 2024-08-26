package nats

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Unlites/ml-analysis-provider/worker/internal/application"
	natsclient "github.com/nats-io/nats.go"
)

const analysisSubject = "analysis"

// AnalysisNats is a NATS representation of Analysis
type AnalysisNats struct {
	Id              string `json:"id"`
	Query           string `json:"query"`
	Answer          string `json:"answer"`
	IsUserSatisfied bool   `json:"is_user_satisfied"`
}

// AnalyzesFilterNats is a NATS representation of AnalysisFilter
type AnalyzesFilterNats struct {
	Query           string `json:"query"`
	Answer          string `json:"answer"`
	IsUserSatisfied *bool  `json:"is_user_satisfied"`
}

// NatsMqHandler is a nats adapter as handler for application
type NatsMqHandler struct {
	handleTimeout time.Duration
	conn          *natsclient.Conn
	uc            application.Usecase
}

// NewNatsMqHandler creates new nats handler and subscribes to required subjects
func NewMqHandlerNats(
	conn *natsclient.Conn,
	handleTimeout time.Duration,
	uc application.Usecase,
) (*NatsMqHandler, error) {
	handler := &NatsMqHandler{
		handleTimeout: handleTimeout,
		conn:          conn,
		uc:            uc,
	}

	conn.Subscribe(analysisSubject, handler.AddAnalysis)
	conn.Subscribe(analysisSubject+".filter", handler.GetAnalyzes)
	conn.Subscribe(analysisSubject+".*", handler.GetAnalysisById)

	return handler, nil
}

// AddAnalysis parses NATS message and calls usecase to add analysis
func (h *NatsMqHandler) AddAnalysis(m *natsclient.Msg) {
	var analysis AnalysisNats

	if err := json.Unmarshal(m.Data, &analysis); err != nil {
		slog.Error("failed to unmarshal json to anylysis", "detail", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.handleTimeout)
	defer cancel()

	if err := h.uc.AddAnalysis(ctx, toDomainAnalysis(analysis)); err != nil {
		slog.Error("failed to add anylysis", "detail", err)
		return
	}
}

// GetAnalysisById gets an analysis by id
func (h *NatsMqHandler) GetAnalysisById(m *natsclient.Msg) {
	id := m.Subject[len(analysisSubject)+1:]

	ctx, cancel := context.WithTimeout(context.Background(), h.handleTimeout)
	defer cancel()

	analysis, err := h.uc.GetAnalysisById(ctx, id)
	if err != nil {
		slog.Error("failed to get analysis", "detail", err)
		return
	}

	analysisResponse, err := json.Marshal(toAnalysisNats(analysis))
	if err != nil {
		if err := m.Respond(nil); err != nil {
			slog.Error("failed to respond error", "detail", err)
		}

		slog.Error("failed to marshal analysis to json", "detail", err)
		return
	}

	if err := m.Respond(analysisResponse); err != nil {
		slog.Error("failed to respond analysis", "detail", err)
	}
}

// GetAnalyzes gets analyzes by filter
func (h *NatsMqHandler) GetAnalyzes(m *natsclient.Msg) {
	var filter AnalyzesFilterNats

	if err := json.Unmarshal(m.Data, &filter); err != nil {
		// TODO: respond error
		if err := m.Respond(nil); err != nil {
			slog.Error("failed to respond error", "detail", err)
		}

		slog.Error("failed to marshal analysis to json", "detail", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.handleTimeout)
	defer cancel()

	analyzes, err := h.uc.GetAnalyzes(ctx, toDomainAnalyzesFilter(filter))
	if err != nil {
		if err := m.Respond(nil); err != nil {
			slog.Error("failed to respond error", "detail", err)
		}

		slog.Error("failed to get analyzes", "detail", err)
		return
	}

	analyzesResponse, err := json.Marshal(toAnalyzesNats(analyzes))
	if err != nil {
		if err := m.Respond(nil); err != nil {
			slog.Error("failed to respond error", "detail", err)
		}

		slog.Error("failed to marshal analyzes", "detail", err)
		return
	}

	if err := m.Respond(analyzesResponse); err != nil {
		slog.Error("failed to respond analyzes", "detail", err)
	}
}
