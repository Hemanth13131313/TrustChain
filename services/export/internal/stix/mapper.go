package stix

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type STIXIncident struct {
	Type        string    `json:"type"` // "incident"
	ID          string    `json:"id"`
	Created     time.Time `json:"created"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Mapper struct {
	pool *pgxpool.Pool
}

func NewMapper(pool *pgxpool.Pool) *Mapper {
	return &Mapper{pool: pool}
}

func (m *Mapper) ExportIncidents(ctx context.Context) ([]STIXIncident, error) {
	rows, err := m.pool.Query(ctx, "SELECT id, event_type, created_at FROM audit_logs ORDER BY created_at DESC LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []STIXIncident
	for rows.Next() {
		var id int
		var eventType string
		var createdAt time.Time
		
		if err := rows.Scan(&id, &eventType, &createdAt); err != nil {
			return nil, err
		}

		incidents = append(incidents, STIXIncident{
			Type:        "incident",
			ID:          fmt.Sprintf("incident--trustchain-%d", id),
			Created:     createdAt,
			Name:        eventType,
			Description: "TRUSTCHAIN Policy Violation / Drift Event",
		})
	}
	return incidents, nil
}
