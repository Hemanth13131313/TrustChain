package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trustchain/verification/internal/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// SaveSignature saves the signature metadata
func (r *PostgresRepository) SaveSignature(ctx context.Context, sig domain.SignatureInfo) error {
	query := `
		INSERT INTO signatures (artifact_id, subject, issuer, verified)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(ctx, query, sig.ArtifactID, sig.Subject, sig.Issuer, sig.Verified)
	if err != nil {
		return fmt.Errorf("inserting signature: %w", err)
	}
	return nil
}

// SaveProvenance saves the provenance attestation metadata
func (r *PostgresRepository) SaveProvenance(ctx context.Context, prov domain.SLSAInfo) error {
	query := `
		INSERT INTO provenance_attestations (artifact_id, builder_id, source_repo, slsa_level, verified)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, prov.ArtifactID, prov.BuilderID, prov.SourceRepo, prov.SLSALevel, prov.Verified)
	if err != nil {
		return fmt.Errorf("inserting provenance: %w", err)
	}
	return nil
}

// GetArtifactIDByDigest fetches the artifact ID given a digest
func (r *PostgresRepository) GetArtifactIDByDigest(ctx context.Context, digest string) (string, error) {
	var id string
	query := `SELECT id FROM artifacts WHERE digest = $1 LIMIT 1`
	err := r.pool.QueryRow(ctx, query, digest).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("getting artifact id by digest: %w", err)
	}
	return id, nil
}
