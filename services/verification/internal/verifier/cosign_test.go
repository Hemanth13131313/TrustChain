package verifier

import (
	"context"
	"testing"
)

func TestVerifySignature_Mock(t *testing.T) {
	v := NewCosignVerifier()

	sigs, err := v.VerifySignature(context.Background(), "mock-image", "test-artifact-id")
	if err != nil {
		t.Fatalf("expected no error for mock-image, got: %v", err)
	}

	if len(sigs) != 1 {
		t.Fatalf("expected 1 signature, got %d", len(sigs))
	}

	sig := sigs[0]
	if sig.Subject != "test@trustchain.dev" {
		t.Errorf("unexpected subject: %s", sig.Subject)
	}
	if sig.ArtifactID != "test-artifact-id" {
		t.Errorf("unexpected artifact id: %s", sig.ArtifactID)
	}
	if !sig.Verified {
		t.Errorf("expected signature to be verified")
	}
}

func TestVerifySignature_Unimplemented(t *testing.T) {
	v := NewCosignVerifier()
	_, err := v.VerifySignature(context.Background(), "real-image:latest", "test-artifact-id")
	if err == nil {
		t.Fatalf("expected error for unimplemented live registry, got nil")
	}
}
