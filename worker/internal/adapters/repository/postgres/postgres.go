package postgres

import (
	"context"
	"fmt"

	"github.com/Unlites/ml-analysis-provider/worker/internal/domain"
	"github.com/jackc/pgx/v5"
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

	query := "INSERT INTO ml_analysis.analyzes (query, answer, is_user_satisfied) VALUES ($1, $2, $3)"

	_, err = conn.Exec(ctx, query, analysis.Query, analysis.Answer, analysis.IsUserSatisfied)
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

	query := "SELECT id, query, answer, is_user_satisfied FROM ml_analysis.analyzes WHERE id = $1"

	err = conn.QueryRow(ctx, query, id).Scan(
		&analysis.Id,
		&analysis.Query,
		&analysis.Answer,
		&analysis.IsUserSatisfied,
	)
	if err != nil {
		return domain.Analysis{}, fmt.Errorf("failed to scan from postgres: %w", err)
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

	query := "SELECT id, query, answer, is_user_satisfied FROM ml_analysis.analyzes WHERE 1=1"
	args := make(pgx.NamedArgs)

	if filter.Query != "" {
		query += " AND query ILIKE @query"
		args["query"] = "%" + filter.Query + "%"
	}
	if filter.Answer != "" {
		query += " AND answer ILIKE @answer"
		args["answer"] = "%" + filter.Answer + "%"
	}
	if filter.IsUserSatisfied != nil {
		query += " AND is_user_satisfied = @is_user_satisfied"
		args["is_user_satisfied"] = *filter.IsUserSatisfied
	}
	if filter.Limit > 0 {
		query += " LIMIT @limit"
		args["limit"] = filter.Limit
	}
	if filter.Offset > 0 {
		query += " OFFSET @offset"
		args["offset"] = filter.Offset
	}

	if filter.IsUserSatisfied != nil {
		args["is_user_satisfied"] = *filter.IsUserSatisfied
	}

	rows, err := conn.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to do query in postgres: %w", err)
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

	query := "SELECT id, query, answer, is_user_satisfied FROM ml_analysis.analyzes WHERE id = ANY($1)"

	rows, err := conn.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to do query in postgres: %w", err)
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
