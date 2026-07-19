package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trustchain/policy-engine/internal/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// GatherContext fetches signature and provenance data for the artifact digest
func (r *PostgresRepository) GatherContext(ctx context.Context, digest string) (*domain.ArtifactContext, error) {
	var artCtx domain.ArtifactContext
	artCtx.Digest = digest

	// Get artifact ID
	var artifactID string
	err := r.pool.QueryRow(ctx, "SELECT id FROM artifacts WHERE digest = $1 LIMIT 1", digest).Scan(&artifactID)
	if err != nil {
		return nil, fmt.Errorf("getting artifact: %w", err)
	}

	// Get signature counts
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM signatures WHERE artifact_id = $1 AND verified = true", artifactID).Scan(&artCtx.SignaturesCount)
	if err != nil {
		return nil, fmt.Errorf("getting signature count: %w", err)
	}
	artCtx.SignatureVerified = artCtx.SignaturesCount > 0

	// Get SLSA level (max verified level)
	err = r.pool.QueryRow(ctx, "SELECT COALESCE(MAX(slsa_level), 0) FROM provenance_attestations WHERE artifact_id = $1 AND verified = true", artifactID).Scan(&artCtx.SLSALevel)
	if err != nil {
		return nil, fmt.Errorf("getting slsa level: %w", err)
	}

	// For MVP, we simulate vulnerability correlation by mocking a VEX response.
	// Real implementation would join sbom_documents -> vulnerabilites.
	if digest == "mock-vulnerable-digest" {
		artCtx.Vulnerabilities = []domain.Vulnerability{
			{CVE: "CVE-2024-1234", Severity: "CRITICAL", Status: "UNMITIGATED"},
		}
	} else {
		artCtx.Vulnerabilities = []domain.Vulnerability{}
	}

	return &artCtx, nil
}
