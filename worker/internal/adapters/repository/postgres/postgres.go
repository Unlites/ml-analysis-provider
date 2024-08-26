package repository

import (
	"context"
	"fmt"

	"github.com/Unlites/ml-analysis-provider/worker/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

// AddAnalysis adds an new analysis from Postgres
func (r *PostgresRepository) AddAnalysis(ctx context.Context, analysis domain.Analysis) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire pg connection: %w", err)
	}
	defer conn.Release()

	query := "INSERT INTO analysis (id, query, answer, is_user_satisfied) VALUES ($1, $2, $3, $4)"

	_, err = conn.Exec(ctx, query, analysis.Id, analysis.Query, analysis.Answer, analysis.IsUserSatisfied)
	if err != nil {
		return fmt.Errorf("failed to add analysis: %w", err)
	}

	return nil
}

// GetAnalysisById gets an analysis by id from Postgres
func (r *PostgresRepository) GetAnalysisById(ctx context.Context, id string) (domain.Analysis, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return domain.Analysis{}, fmt.Errorf("failed to acquire pg connection: %w", err)
	}
	defer conn.Release()

	var analysis domain.Analysis

	query := "SELECT id, query, answer, is_user_satisfied FROM analysis WHERE id = $1"

	err = conn.QueryRow(ctx, query, id).Scan(
		&analysis.Id,
		&analysis.Query,
		&analysis.Answer,
		&analysis.IsUserSatisfied,
	)
	if err != nil {
		return domain.Analysis{}, fmt.Errorf("failed to get analysis by id: %w", err)
	}

	return analysis, nil
}

// GetAnalyzes gets analyzes by filter from PostgreSQL
func (r *PostgresRepository) GetAnalyzes(
	ctx context.Context,
	filter domain.AnalyzesFilter,
) ([]domain.Analysis, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire pg connection: %w", err)
	}
	defer conn.Release()

	var analyzes []domain.Analysis

	query := "SELECT id, query, answer, is_user_satisfied FROM analysis"

	if filter.Query != "" {
		query += " WHERE query ILIKE $1"
	}
	if filter.Answer != "" {
		query += " WHERE answer ILIKE $2"
	}
	if filter.IsUserSatisfied != nil {
		query += " WHERE is_user_satisfied = $3"
	}
	if filter.Limit > 0 {
		query += " LIMIT $4"
	}
	if filter.Offset > 0 {
		query += " OFFSET $5"
	}

	rows, err := conn.Query(ctx, query,
		"%"+filter.Query+"%",
		"%"+filter.Answer+"%",
		filter.IsUserSatisfied,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get analyzes: %w", err)
	}

	for rows.Next() {
		var analysis domain.Analysis
		err = rows.Scan(
			&analysis.Id,
			&analysis.Query,
			&analysis.Answer,
			&analysis.IsUserSatisfied,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan analysis: %w", err)
		}

		analyzes = append(analyzes, analysis)
	}

	return analyzes, nil
}

// GetAnalyzesByIds gets analyzes by given slices of ids from PostgreSQL
func (r *PostgresRepository) GetAnalyzesByIds(ctx context.Context, ids []string) ([]domain.Analysis, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire pg connection: %w", err)
	}
	defer conn.Release()

	var analyzes []domain.Analysis

	query := "SELECT id, query, answer, is_user_satisfied FROM analysis WHERE id = ANY($1)"

	rows, err := conn.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get analyzes by ids: %w", err)
	}

	for rows.Next() {
		var analysis domain.Analysis
		err = rows.Scan(
			&analysis.Id,
			&analysis.Query,
			&analysis.Answer,
			&analysis.IsUserSatisfied,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan analysis: %w", err)
		}

		analyzes = append(analyzes, analysis)
	}

	return analyzes, nil
}
