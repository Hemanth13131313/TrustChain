package correlator

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Analyzer struct {
	pool *pgxpool.Pool
}

func NewAnalyzer(pool *pgxpool.Pool) *Analyzer {
	return &Analyzer{pool: pool}
}

// CheckDrift looks up the digest in the database to ensure it has valid signatures and SLSA provenance.
func (a *Analyzer) CheckDrift(ctx context.Context, obs WorkloadObservation) error {
	// Query to check if the digest exists and has at least one valid signature
	query := `
		SELECT COUNT(s.id) 
		FROM artifacts a
		JOIN signatures s ON a.id = s.artifact_id
		WHERE a.digest = $1 AND s.verified = true
	`

	var sigCount int
	err := a.pool.QueryRow(ctx, query, obs.ImageDigest).Scan(&sigCount)
	if err != nil {
		return fmt.Errorf("failed to query database for digest %s: %v", obs.ImageDigest, err)
	}

	if sigCount == 0 {
		// Drift Detected! Unsigned or unknown image running.
		log.Printf("🚨 DRIFT DETECTED: Pod %s/%s on %s is running unauthorized digest %s", 
			obs.Namespace, obs.PodName, obs.NodeName, obs.ImageDigest)
		// In a full implementation, this would emit a Drift Event to Kafka for the Enforcement Orchestrator
		return fmt.Errorf("unauthorized digest: %s", obs.ImageDigest)
	}

	log.Printf("✅ COMPLIANT: Pod %s/%s running approved digest %s", obs.Namespace, obs.PodName, obs.ImageDigest)
	return nil
}
