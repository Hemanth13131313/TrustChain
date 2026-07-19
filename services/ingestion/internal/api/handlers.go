package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/trustchain/ingestion/internal/domain"
	"github.com/trustchain/ingestion/internal/parser"
	"github.com/trustchain/ingestion/internal/storage"
)

type IngestionHandler struct {
	db   *storage.PostgresRepository
	blob storage.BlobRepository
}

func NewIngestionHandler(db *storage.PostgresRepository, blob storage.BlobRepository) *IngestionHandler {
	return &IngestionHandler{db: db, blob: blob}
}

func (h *IngestionHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/artifacts/{digest}/sbom", h.HandleIngestSBOM)
	r.Post("/api/v1/artifacts/{digest}/provenance", h.HandleIngestProvenance)
}

func (h *IngestionHandler) HandleIngestSBOM(w http.ResponseWriter, r *http.Request) {
	digest := chi.URLParam(r, "digest")
	registry := r.URL.Query().Get("registry")
	repository := r.URL.Query().Get("repository")
	format := r.URL.Query().Get("format") // cyclonedx, spdx-json, spdx-tv

	if registry == "" || repository == "" || format == "" {
		http.Error(w, "missing required query parameters: registry, repository, format", http.StatusBadRequest)
		return
	}

	// 1. Upsert Artifact
	art := domain.Artifact{
		Digest:     digest,
		Registry:   registry,
		Repository: repository,
	}
	artID, err := h.db.UpsertArtifact(r.Context(), art)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to upsert artifact: %v", err), http.StatusInternalServerError)
		return
	}

	// 2. Read request body into buffer so we can parse and store
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	// 3. Parse SBOM
	var components []domain.NormalizedComponent
	if format == "cyclonedx" {
		components, err = parser.ParseCycloneDX(bytes.NewReader(body))
	} else if format == "spdx-json" {
		components, err = parser.ParseSPDX(bytes.NewReader(body), "json")
	} else if format == "spdx-tv" {
		components, err = parser.ParseSPDX(bytes.NewReader(body), "tag-value")
	} else {
		http.Error(w, "unsupported format", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse SBOM: %v", err), http.StatusBadRequest)
		return
	}

	// 4. Save raw document to blob storage
	objectName := fmt.Sprintf("sboms/%s/%s-%d", registry, digest, time.Now().Unix())
	storageRef, err := h.blob.Save(r.Context(), objectName, bytes.NewReader(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to save blob: %v", err), http.StatusInternalServerError)
		return
	}

	// 5. Save normalized metadata to Postgres
	doc := domain.SBOMDocument{
		ArtifactID:           artID,
		Format:               format,
		NormalizedComponents: components,
		StorageRef:           storageRef,
	}
	if err := h.db.SaveSBOMDocument(r.Context(), doc); err != nil {
		http.Error(w, fmt.Sprintf("failed to save metadata: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":      "success",
		"artifact_id": artID,
		"storage_ref": storageRef,
	})
}

func (h *IngestionHandler) HandleIngestProvenance(w http.ResponseWriter, r *http.Request) {
	// MVP: stub for provenance ingestion
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error": "not implemented yet"}`))
}
