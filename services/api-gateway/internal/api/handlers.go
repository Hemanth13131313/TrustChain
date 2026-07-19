package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/trustchain/api-gateway/internal/auth"
	"github.com/trustchain/api-gateway/internal/storage"
)

type DashboardHandler struct {
	db   *storage.PostgresRepository
	auth *auth.JWTAuth
}

func NewDashboardHandler(db *storage.PostgresRepository, auth *auth.JWTAuth) *DashboardHandler {
	return &DashboardHandler{db: db, auth: auth}
}

func (h *DashboardHandler) RegisterRoutes(r chi.Router) {
	// Public stub for dev
	r.Post("/api/v1/auth/login", h.HandleLogin)

	// Protected endpoints
	r.Group(func(r chi.Router) {
		r.Use(h.auth.Middleware)
		r.Get("/api/v1/dashboard/metrics", h.HandleGetMetrics)
		r.Get("/api/v1/artifacts", h.HandleListArtifacts)
	})
}

// HandleLogin is a stub to generate a token for local dev
func (h *DashboardHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	token, err := h.auth.GenerateToken("admin")
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *DashboardHandler) HandleGetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.db.GetDashboardMetrics(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (h *DashboardHandler) HandleListArtifacts(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	artifacts, err := h.db.ListArtifacts(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "failed to fetch artifacts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"artifacts": artifacts,
		"limit":     limit,
		"offset":    offset,
	})
}
