package domain

import "time"

// Artifact represents a minimal version of Artifact needed for verification
type Artifact struct {
	ID         string
	Digest     string
	Registry   string
	Repository string
}

// VerificationResult represents the payload emitted when verification completes
type VerificationResult struct {
	ArtifactID      string    `json:"artifact_id"`
	Digest          string    `json:"digest"`
	Verified        bool      `json:"verified"`
	SignaturesCount int       `json:"signatures_count"`
	SLSALevel       int       `json:"slsa_level"`
	Reason          string    `json:"reason,omitempty"`
	VerifiedAt      time.Time `json:"verified_at"`
}

// SignatureInfo represents metadata extracted from a verified signature
type SignatureInfo struct {
	ArtifactID string
	Subject    string
	Issuer     string
	Verified   bool
}

// SLSAInfo represents metadata extracted from a verified provenance attestation
type SLSAInfo struct {
	ArtifactID string
	BuilderID  string
	SourceRepo string
	SLSALevel  int
	Verified   bool
}
