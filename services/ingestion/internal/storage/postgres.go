package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trustchain/ingestion/internal/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// UpsertArtifact inserts or returns the existing artifact ID for a given digest/registry/repo
func (r *PostgresRepository) UpsertArtifact(ctx context.Context, art domain.Artifact) (string, error) {
	var id string
	query := `
		INSERT INTO artifacts (digest, registry, repository)
		VALUES ($1, $2, $3)
		ON CONFLICT (digest, registry, repository) DO UPDATE SET first_seen = artifacts.first_seen
		RETURNING id
	`
	err := r.pool.QueryRow(ctx, query, art.Digest, art.Registry, art.Repository).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("upserting artifact: %w", err)
	}
	return id, nil
}

// SaveSBOMDocument saves the normalized SBOM metadata
func (r *PostgresRepository) SaveSBOMDocument(ctx context.Context, doc domain.SBOMDocument) error {
	componentsJSON, err := json.Marshal(doc.NormalizedComponents)
	if err != nil {
		return fmt.Errorf("marshaling components: %w", err)
	}

	query := `
		INSERT INTO sbom_documents (artifact_id, format, normalized_components, storage_ref)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	err = r.pool.QueryRow(ctx, query, doc.ArtifactID, doc.Format, componentsJSON, doc.StorageRef).Scan(&id)
	if err != nil {
		return fmt.Errorf("inserting sbom document: %w", err)
	}
	return nil
}

// SaveProvenance saves the provenance attestation metadata
func (r *PostgresRepository) SaveProvenance(ctx context.Context, prov domain.ProvenanceAttestation) error {
	query := `
		INSERT INTO provenance_attestations (artifact_id, builder_id, source_repo, slsa_level, verified)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id string
	err := r.pool.QueryRow(ctx, query, prov.ArtifactID, prov.BuilderID, prov.SourceRepo, prov.SLSALevel, prov.Verified).Scan(&id)
	if err != nil {
		return fmt.Errorf("inserting provenance: %w", err)
	}
	return nil
}
