package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Unlites/ml-analysis-provider/worker/internal/domain"
)

// Usecase is a business logic of the application
type Usecase interface {
	// AddAnalysis adds an new analysis
	AddAnalysis(ctx context.Context, analysis domain.Analysis) error

	// GetAnalysisById gets an analysis by id
	GetAnalysisById(ctx context.Context, id string) (domain.Analysis, error)

	// GetAnalyzes gets analyzes by filter
	GetAnalyzes(ctx context.Context, filter domain.AnalyzesFilter) ([]domain.Analysis, error)
}

// FullTextSearcher searches analyzes by filter using full-text search
type FullTextSearcher interface {
	// SearchAnalyzes searches analyzes by given filter and returnes slice of founded analyzes ids
	SearchAnalyzes(ctx context.Context, filter domain.AnalyzesFilter) (ids []string, err error)
}

// Repository works with analyzes storage
type Repository interface {
	// AddAnalysis adds an new analysis
	AddAnalysis(ctx context.Context, analysis domain.Analysis) error

	// GetAnalysisById gets an analysis by id
	GetAnalysisById(ctx context.Context, id string) (domain.Analysis, error)

	// GetAnalyzes gets analyzes by filter
	GetAnalyzes(ctx context.Context, filter domain.AnalyzesFilter) ([]domain.Analysis, error)

	// GetAnalyzesByIds gets analyzes by given slices of ids
	GetAnalyzesByIds(ctx context.Context, ids []string) ([]domain.Analysis, error)
}

type usecase struct {
	searcher FullTextSearcher
	repo     Repository
}

// NewUsecase creates new usecase
func NewUsecase(searcher FullTextSearcher, repo Repository) *usecase {
	return &usecase{
		searcher: searcher,
		repo:     repo,
	}
}

// AddAnalysis implements Usecase
func (u *usecase) AddAnalysis(ctx context.Context, analysis domain.Analysis) error {
	return u.repo.AddAnalysis(ctx, analysis)
}

// GetAnalysisById implements Usecase
func (u *usecase) GetAnalysisById(ctx context.Context, id string) (domain.Analysis, error) {
	return u.repo.GetAnalysisById(ctx, id)
}

// GetAnalyzes implements Usecase
func (u *usecase) GetAnalyzes(ctx context.Context, filter domain.AnalyzesFilter) ([]domain.Analysis, error) {
	var analyzes []domain.Analysis

	ids, err := u.searcher.SearchAnalyzes(ctx, filter)
	if err != nil || len(ids) == 0 {
		slog.Error("failed to search analyzes", "detail", err)

		analyzes, err = u.repo.GetAnalyzes(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to get analyzes from repo: %w", err)
		}
	} else {
		analyzes, err = u.repo.GetAnalyzesByIds(ctx, ids)
		if err != nil {
			return nil, fmt.Errorf("failed to get analyzes by ids from repo: %w", err)
		}
	}

	return analyzes, nil
}
