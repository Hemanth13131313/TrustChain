package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trustchain/api-gateway/internal/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// GetDashboardMetrics aggregates global counts
func (r *PostgresRepository) GetDashboardMetrics(ctx context.Context) (*domain.DashboardMetrics, error) {
	var metrics domain.DashboardMetrics

	// Count total artifacts
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM artifacts").Scan(&metrics.TotalArtifacts)
	if err != nil {
		return nil, fmt.Errorf("counting artifacts: %w", err)
	}

	// Count total signatures
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM signatures WHERE verified = true").Scan(&metrics.TotalSignatures)
	if err != nil {
		return nil, fmt.Errorf("counting signatures: %w", err)
	}

	// Count total verified artifacts (artifacts with at least one verified signature)
	queryVerified := `
		SELECT COUNT(DISTINCT artifact_id) 
		FROM signatures 
		WHERE verified = true
	`
	err = r.pool.QueryRow(ctx, queryVerified).Scan(&metrics.TotalVerified)
	if err != nil {
		return nil, fmt.Errorf("counting verified artifacts: %w", err)
	}

	return &metrics, nil
}

// ListArtifacts retrieves the artifacts with their aggregated status
func (r *PostgresRepository) ListArtifacts(ctx context.Context, limit, offset int) ([]domain.ArtifactView, error) {
	query := `
		SELECT 
			a.id, a.digest, a.registry, a.repository, a.first_seen,
			COALESCE(s.sig_count, 0) as signatures_count,
			COALESCE(p.max_slsa, 0) as slsa_level,
			CASE WHEN sb.id IS NOT NULL THEN true ELSE false END as has_sbom
		FROM artifacts a
		LEFT JOIN (
			SELECT artifact_id, COUNT(*) as sig_count FROM signatures WHERE verified = true GROUP BY artifact_id
		) s ON a.id = s.artifact_id
		LEFT JOIN (
			SELECT artifact_id, MAX(slsa_level) as max_slsa FROM provenance_attestations WHERE verified = true GROUP BY artifact_id
		) p ON a.id = p.artifact_id
		LEFT JOIN (
			SELECT DISTINCT artifact_id, id FROM sbom_documents
		) sb ON a.id = sb.artifact_id
		ORDER BY a.first_seen DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("querying artifacts: %w", err)
	}
	defer rows.Close()

	var results []domain.ArtifactView
	for rows.Next() {
		var view domain.ArtifactView
		err := rows.Scan(
			&view.ID, &view.Digest, &view.Registry, &view.Repository, &view.FirstSeen,
			&view.SignaturesCount, &view.SLSALevel, &view.HasSBOM,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning artifact row: %w", err)
		}
		results = append(results, view)
	}

	return results, nil
}
