package orchestrator

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Auditor struct {
	pool *pgxpool.Pool
}

func NewAuditor(pool *pgxpool.Pool) *Auditor {
	return &Auditor{pool: pool}
}

// LogIncident writes a tamper-evident log entry.
func (a *Auditor) LogIncident(ctx context.Context, eventType string, payload map[string]interface{}) error {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Get the last hash
	var lastHash string
	err = tx.QueryRow(ctx, "SELECT hash FROM audit_logs ORDER BY id DESC LIMIT 1 FOR UPDATE").Scan(&lastHash)
	if err != nil {
		return fmt.Errorf("failed to get previous hash: %v", err)
	}

	payloadBytes, _ := json.Marshal(payload)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Compute new hash: SHA256(lastHash + eventType + payload + timestamp)
	dataToHash := fmt.Sprintf("%s|%s|%s|%s", lastHash, eventType, string(payloadBytes), timestamp)
	hash := sha256.Sum256([]byte(dataToHash))
	hashStr := hex.EncodeToString(hash[:])

	// Insert new log
	_, err = tx.Exec(ctx, 
		"INSERT INTO audit_logs (event_type, payload, previous_hash, hash, created_at) VALUES ($1, $2, $3, $4, $5)",
		eventType, string(payloadBytes), lastHash, hashStr, timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert audit log: %v", err)
	}

	return tx.Commit(ctx)
}
