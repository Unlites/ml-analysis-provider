package application

import (
	"context"

	"github.com/Unlites/ml-analysis-provider/controller/internal/domain"
)

// Usecase is the business logic of the application
type Usecase interface {
	// AddAnalysis adds an new analysis
	AddAnalysis(ctx context.Context, analysis domain.Analysis) error

	// GetAnalysisById gets an analysis by id
	GetAnalysisById(ctx context.Context, id string) (domain.Analysis, error)

	// GetAnalyzes gets analyzes by filter
	GetAnalyzes(ctx context.Context, AnalyzesFilter domain.AnalyzesFilter) ([]domain.Analysis, error)
}

// AnalysisBroker is the broker that publishes and requests analyzes
type AnalysisBroker interface {
	// PublishAnalysis publishes an analysis to the broker
	PublishAnalysis(ctx context.Context, analysis domain.Analysis) error

	// RequestAnalysisById requests an analysis by id
	RequestAnalysisById(ctx context.Context, id string) (domain.Analysis, error)

	// RequestAnalyzes requests analyzes by filter
	RequestAnalyzes(ctx context.Context, AnalyzesFilter domain.AnalyzesFilter) ([]domain.Analysis, error)
}

type usecase struct {
	broker AnalysisBroker
}

// AddAnalysis implements Usecase
func (u *usecase) AddAnalysis(ctx context.Context, analysis domain.Analysis) error {
	return u.broker.PublishAnalysis(ctx, analysis)
}

// GetAnalysisById implements Usecase
func (u *usecase) GetAnalysisById(ctx context.Context, id string) (domain.Analysis, error) {
	return u.broker.RequestAnalysisById(ctx, id)
}

// GetAnalyzes implements Usecase
func (u *usecase) GetAnalyzes(ctx context.Context, AnalyzesFilter domain.AnalyzesFilter) ([]domain.Analysis, error) {
	return u.broker.RequestAnalyzes(ctx, AnalyzesFilter)
}

// NewUsecase creates a new usecase object representing the business logic of the application
func NewUsecase(broker AnalysisBroker) *usecase {
	return &usecase{
		broker: broker,
	}
}
