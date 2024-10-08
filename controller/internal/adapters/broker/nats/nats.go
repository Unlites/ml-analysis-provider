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

// analysisResponse is a NATS response for analysis request
type analysisResponse struct {
	Analysis AnalysisNats `json:"data"`
	Error    string       `json:"error"`
}

// analysisResponse is a NATS response for multiple analyzes request
type analyzesResponse struct {
	Analyzes []AnalysisNats `json:"data"`
	Error    string         `json:"error"`
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
	analysisMsg, err := n.conn.RequestWithContext(ctx, analysisSubject+".id."+id, nil)
	if err != nil {
		return domain.Analysis{}, fmt.Errorf("failed to do request: %w", err)
	}

	var response analysisResponse

	if err := json.Unmarshal(analysisMsg.Data, &response); err != nil {
		return domain.Analysis{}, fmt.Errorf("failed to unmarshal analysis from json: %w", err)
	}

	if response.Error != "" {
		return domain.Analysis{}, fmt.Errorf("error in response: %v", response.Error)
	}

	return toDomainAnalysis(response.Analysis), nil
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

	var response analyzesResponse
	if err := json.Unmarshal(analysisMsg.Data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("error in response: %v", response.Error)
	}

	return toDomainAnalyzes(response.Analyzes), nil
}
