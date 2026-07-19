package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/trustchain/verification/internal/domain"
	"github.com/trustchain/verification/internal/event"
	"github.com/trustchain/verification/internal/storage"
	"github.com/trustchain/verification/internal/verifier"
)

type VerificationHandler struct {
	db        *storage.PostgresRepository
	publisher *event.KafkaPublisher
	cosign    *verifier.CosignVerifier
	slsa      *verifier.SLSAVerifier
}

func NewVerificationHandler(
	db *storage.PostgresRepository,
	publisher *event.KafkaPublisher,
	cosign *verifier.CosignVerifier,
	slsa *verifier.SLSAVerifier,
) *VerificationHandler {
	return &VerificationHandler{
		db:        db,
		publisher: publisher,
		cosign:    cosign,
		slsa:      slsa,
	}
}

func (h *VerificationHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/artifacts/{digest}/verify", h.HandleVerify)
}

type VerifyRequest struct {
	ImageRef string `json:"image_ref"`
}

func (h *VerificationHandler) HandleVerify(w http.ResponseWriter, r *http.Request) {
	digest := chi.URLParam(r, "digest")

	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 1. Look up artifact ID
	artifactID, err := h.db.GetArtifactIDByDigest(r.Context(), digest)
	if err != nil {
		http.Error(w, fmt.Sprintf("artifact not found for digest %s", digest), http.StatusNotFound)
		return
	}

	// 2. Verify Signatures
	sigs, err := h.cosign.VerifySignature(r.Context(), req.ImageRef, artifactID)
	if err != nil {
		http.Error(w, fmt.Sprintf("verifying signature: %v", err), http.StatusInternalServerError)
		return
	}
	
	for _, sig := range sigs {
		_ = h.db.SaveSignature(r.Context(), sig) // ignoring error for brevity in scaffold
	}

	// 3. Verify Provenance
	prov, err := h.slsa.VerifyProvenance(r.Context(), req.ImageRef, artifactID)
	if err != nil {
		http.Error(w, fmt.Sprintf("verifying provenance: %v", err), http.StatusInternalServerError)
		return
	}
	_ = h.db.SaveProvenance(r.Context(), *prov)

	// 4. Determine overall result
	verified := len(sigs) > 0 && prov != nil && prov.Verified

	result := domain.VerificationResult{
		ArtifactID:      artifactID,
		Digest:          digest,
		Verified:        verified,
		SignaturesCount: len(sigs),
		SLSALevel:       prov.SLSALevel,
		VerifiedAt:      time.Now(),
	}

	// 5. Publish Event
	if err := h.publisher.PublishVerificationResult(r.Context(), result); err != nil {
		http.Error(w, fmt.Sprintf("publishing result: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
