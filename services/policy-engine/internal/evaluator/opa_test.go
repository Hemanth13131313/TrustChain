package evaluator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestOPAEvaluator(t *testing.T) {
	// Create a temporary Rego file for testing
	regoContent := `
package trustchain

default allow = false
default reason = "no policy evaluated to true"

allow {
    input.slsa_level >= 3
}

reason = "insufficient SLSA level" {
    not allow
} else = "allowed" {
    allow
}
`
	tmpDir := t.TempDir()
	policyPath := filepath.Join(tmpDir, "rules.rego")
	if err := os.WriteFile(policyPath, []byte(regoContent), 0644); err != nil {
		t.Fatalf("failed to write mock policy: %v", err)
	}

	evaluator, err := NewOPAEvaluator(policyPath)
	if err != nil {
		t.Fatalf("failed to create evaluator: %v", err)
	}

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectAllow  bool
		expectReason string
	}{
		{
			name:         "Allowed - SLSA 3",
			input:        map[string]interface{}{"slsa_level": 3},
			expectAllow:  true,
			expectReason: "allowed",
		},
		{
			name:         "Denied - SLSA 2",
			input:        map[string]interface{}{"slsa_level": 2},
			expectAllow:  false,
			expectReason: "insufficient SLSA level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, reason, err := evaluator.Evaluate(context.Background(), tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if allowed != tt.expectAllow {
				t.Errorf("expected allowed %v, got %v", tt.expectAllow, allowed)
			}
			if reason != tt.expectReason {
				t.Errorf("expected reason %q, got %q", tt.expectReason, reason)
			}
		})
	}
}
