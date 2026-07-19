package domain

import "time"

// ArtifactView represents the aggregated view of a single artifact for the dashboard
type ArtifactView struct {
	ID              string    `json:"id"`
	Digest          string    `json:"digest"`
	Registry        string    `json:"registry"`
	Repository      string    `json:"repository"`
	FirstSeen       time.Time `json:"first_seen"`
	SignaturesCount int       `json:"signatures_count"`
	SLSALevel       int       `json:"slsa_level"`
	HasSBOM         bool      `json:"has_sbom"`
}

// DashboardMetrics represents the global aggregated metrics
type DashboardMetrics struct {
	TotalArtifacts  int `json:"total_artifacts"`
	TotalVerified   int `json:"total_verified"`
	TotalSignatures int `json:"total_signatures"`
	// Additional metrics can be added here (e.g. CriticalCVEs)
}
