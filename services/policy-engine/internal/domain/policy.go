package domain

// PolicyRequest represents an incoming request to evaluate policy for an artifact
type PolicyRequest struct {
	Digest string `json:"digest"`
}

// PolicyResponse represents the outcome of the policy evaluation
type PolicyResponse struct {
	Digest  string `json:"digest"`
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
	Cached  bool   `json:"cached"`
}

// ArtifactContext is the aggregated data structure fed into OPA as input
type ArtifactContext struct {
	Digest            string          `json:"digest"`
	SignaturesCount   int             `json:"signatures_count"`
	SignatureVerified bool            `json:"signature_verified"`
	SLSALevel         int             `json:"slsa_level"`
	Vulnerabilities   []Vulnerability `json:"vulnerabilities"`
}

// Vulnerability simulates a VEX entry for an artifact component
type Vulnerability struct {
	CVE      string `json:"cve"`
	Severity string `json:"severity"` // LOW, MEDIUM, HIGH, CRITICAL
	Status   string `json:"status"`   // UNMITIGATED, MITIGATED, FALSE_POSITIVE
}
