package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Unlites/ml-analysis-provider/controller/internal/domain"
	"github.com/nats-io/nats.go"
)

const analysisSubject = "analysis"

// NatsBroker is the broker adapter for NATS
type NatsBroker struct {
	conn *nats.Conn
}

// NewNatsBroker takes a NATS connection and creates a new broker object
func NewNatsBroker(conn *nats.Conn) *NatsBroker {
	return &NatsBroker{
		conn: conn,
	}
}

// PublishAnalysis publishes analysis to NATS
func (n *NatsBroker) PublishAnalysis(ctx context.Context, analysis domain.Analysis) error {
	analysisPayload, err := json.Marshal(analysis)
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

	var analysis domain.Analysis
	if err := json.Unmarshal(analysisMsg.Data, &analysis); err != nil {
		return domain.Analysis{}, fmt.Errorf("failed to unmarshal analysis from json: %w", err)
	}
	return analysis, nil
}

// RequestAnalyzes marshals AnalyzesFilter and sends request to NATS waiting for response with analyzes
func (n *NatsBroker) RequestAnalyzes(ctx context.Context, AnalyzesFilter domain.AnalyzesFilter) ([]domain.Analysis, error) {
	filterPayload, err := json.Marshal(AnalyzesFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal analyzes filter to json: %w", err)
	}

	analysisMsg, err := n.conn.RequestWithContext(ctx, analysisSubject, filterPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	var analysis []domain.Analysis
	if err := json.Unmarshal(analysisMsg.Data, &analysis); err != nil {
		return nil, fmt.Errorf("failed to unmarshal analysis from json: %w", err)
	}
	return analysis, nil
}
