package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Unlites/ml-analysis-provider/controller/internal/domain"
	natsclient "github.com/nats-io/nats.go"
)

const analysisSubject = "analysis"

// NatsBroker is the broker adapter for NATS
type NatsBroker struct {
	conn *natsclient.Conn
}

// NewNatsBroker takes a NATS connection and creates a new broker object
func NewNatsBroker(conn *natsclient.Conn) *NatsBroker {
	return &NatsBroker{
		conn: conn,
	}
}

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

// PublishAnalysis publishes analysis to NATS
func (n *NatsBroker) PublishAnalysis(ctx context.Context, analysis domain.Analysis) error {
	analysisPayload, err := json.Marshal(toAnalysisNats(analysis))
	if err != nil {
		return fmt.Errorf("failed to marshal analysis to json: %w", err)
	}

	return n.conn.Publish(analysisSubject, analysisPayload)
}

// RequestAnalysisById requests NATS for response with particular analysis
func (n *NatsBroker) RequestAnalysisById(ctx context.Context, id string) (domain.Analysis, error) {
	analysisMsg, err := n.conn.RequestWithContext(ctx, analysisSubject+"."+id, nil)
	if err != nil {
		return domain.Analysis{}, fmt.Errorf("failed to do request: %w", err)
	}

	var analysis AnalysisNats

	if err := json.Unmarshal(analysisMsg.Data, &analysis); err != nil {
		return domain.Analysis{}, fmt.Errorf("failed to unmarshal analysis from json: %w", err)
	}
	return toDomainAnalysis(analysis), nil
}

// RequestAnalyzes marshals AnalyzesFilter and sends request to NATS waiting for response with analyzes
func (n *NatsBroker) RequestAnalyzes(ctx context.Context, filter domain.AnalyzesFilter) ([]domain.Analysis, error) {
	filterPayload, err := json.Marshal(toAnalyzesFilterNats(filter))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal analyzes filter to json: %w", err)
	}

	analysisMsg, err := n.conn.RequestWithContext(ctx, analysisSubject+".filter", filterPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	var analysis []AnalysisNats
	if err := json.Unmarshal(analysisMsg.Data, &analysis); err != nil {
		return nil, fmt.Errorf("failed to unmarshal analysis from json: %w", err)
	}
	return toAnalyzesDomain(analysis), nil
}
