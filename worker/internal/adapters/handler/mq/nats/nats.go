package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"unicode/utf8"

	"github.com/Unlites/ml-analysis-provider/worker/internal/application"
	natsclient "github.com/nats-io/nats.go"
)

const analysisSubject = "analysis"

// AnalysisNats is a NATS representation of Analysis
type AnalysisNats struct {
	Id              int    `json:"id"`
	Query           string `json:"query"`
	Answer          string `json:"answer"`
	IsUserSatisfied bool   `json:"is_user_satisfied"`
}

// AnalyzesFilterNats is a NATS representation of AnalysisFilter
type AnalyzesFilterNats struct {
	Query           string `json:"query"`
	Answer          string `json:"answer"`
	IsUserSatisfied *bool  `json:"is_user_satisfied"`
	Limit           int    `json:"limit"`
	Offset          int    `json:"offset"`
}

// Response is a NATS response
type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

// NatsMqHandler is a nats adapter as handler for application
type NatsMqHandler struct {
	handleTimeout time.Duration
	conn          *natsclient.Conn
	uc            application.Usecase
	subs          map[string]*natsclient.Subscription
}

// NewNatsMqHandler creates new nats handler and subscribes to required subjects
func NewNatsMqHandler(
	conn *natsclient.Conn,
	handleTimeout time.Duration,
	uc application.Usecase,
) *NatsMqHandler {

	return &NatsMqHandler{
		handleTimeout: handleTimeout,
		conn:          conn,
		uc:            uc,
		subs:          make(map[string]*natsclient.Subscription),
	}
}

func (h *NatsMqHandler) Start() error {
	subHandlers := map[string]func(*natsclient.Msg){
		analysisSubject:             h.AddAnalysis,
		analysisSubject + ".id.*":   h.GetAnalysisById,
		analysisSubject + ".filter": h.GetAnalyzes,
	}

	for subject, handler := range subHandlers {
		sub, err := h.conn.QueueSubscribe(subject, analysisSubject, handler)
		if err != nil {
			return fmt.Errorf("failed to subscribe to %s subject: %w", subject, err)
		}

		h.subs[subject] = sub
	}

	return nil
}

func (h *NatsMqHandler) Stop() error {
	var err error

	for subject, subscription := range h.subs {
		if err := subscription.Unsubscribe(); err != nil {
			errors.Join(err, fmt.Errorf("failed to unsubscribe from %s: %w", subject, err))
		}
	}

	return err
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
		slog.Error("failed to add analysis: %w", "detail", err)
		return
	}
}

// GetAnalysisById gets an analysis by id
func (h *NatsMqHandler) GetAnalysisById(m *natsclient.Msg) {
	id := m.Subject[utf8.RuneCountInString(analysisSubject+".id."):]

	ctx, cancel := context.WithTimeout(context.Background(), h.handleTimeout)
	defer cancel()

	analysis, err := h.uc.GetAnalysisById(ctx, id)
	if err != nil {
		sendResponse(m, Response{Error: "failed to get analysis: " + err.Error()})
		return
	}

	sendResponse(m, Response{Data: toAnalysisNats(analysis)})
}

// GetAnalyzes gets analyzes by filter
func (h *NatsMqHandler) GetAnalyzes(m *natsclient.Msg) {
	var filter AnalyzesFilterNats

	if err := json.Unmarshal(m.Data, &filter); err != nil {
		sendResponse(m, Response{
			Error: "failed to unmarshal json to anylysis filter: " + err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.handleTimeout)
	defer cancel()

	analyzes, err := h.uc.GetAnalyzes(ctx, toDomainAnalyzesFilter(filter))
	if err != nil {
		sendResponse(m, Response{Error: "failed to get analyzes: " + err.Error()})
		return
	}

	sendResponse(m, Response{Data: toAnalyzesNats(analyzes)})
}
