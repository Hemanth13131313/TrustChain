package domain

import "time"

// Artifact represents a container image or other verifiable artifact
type Artifact struct {
	ID         string    `json:"id"`
	Digest     string    `json:"digest"`
	Registry   string    `json:"registry"`
	Repository string    `json:"repository"`
	TenantID   string    `json:"tenant_id,omitempty"`
	FirstSeen  time.Time `json:"first_seen"`
}

// NormalizedComponent represents a software package extracted from an SBOM
type NormalizedComponent struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	PURL    string `json:"purl"` // Package URL (e.g., pkg:npm/lodash@4.17.21)
	License string `json:"license,omitempty"`
}

// SBOMDocument represents the metadata and storage reference of an ingested SBOM
type SBOMDocument struct {
	ID                   string                `json:"id"`
	ArtifactID           string                `json:"artifact_id"`
	Format               string                `json:"format"` // CycloneDX, SPDX
	NormalizedComponents []NormalizedComponent `json:"normalized_components"`
	StorageRef           string                `json:"storage_ref"`
	CreatedAt            time.Time             `json:"created_at"`
}

// ProvenanceAttestation represents an SLSA provenance document
type ProvenanceAttestation struct {
	ID         string    `json:"id"`
	ArtifactID string    `json:"artifact_id"`
	BuilderID  string    `json:"builder_id"`
	SourceRepo string    `json:"source_repo"`
	SLSALevel  int       `json:"slsa_level"`
	Verified   bool      `json:"verified"`
	CreatedAt  time.Time `json:"created_at"`
}
