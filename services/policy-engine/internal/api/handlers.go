package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/trustchain/policy-engine/internal/domain"
	"github.com/trustchain/policy-engine/internal/evaluator"
	"github.com/trustchain/policy-engine/internal/storage"
)

type PolicyHandler struct {
	db    *storage.PostgresRepository
	cache *storage.RedisCache
	opa   *evaluator.OPAEvaluator
}

func NewPolicyHandler(db *storage.PostgresRepository, cache *storage.RedisCache, opa *evaluator.OPAEvaluator) *PolicyHandler {
	return &PolicyHandler{db: db, cache: cache, opa: opa}
}

func (h *PolicyHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/policy/evaluate", h.HandleEvaluate)
}

func (h *PolicyHandler) HandleEvaluate(w http.ResponseWriter, r *http.Request) {
	var req domain.PolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Digest == "" {
		http.Error(w, "digest is required", http.StatusBadRequest)
		return
	}

	// 1. Check Cache
	if cachedResp, err := h.cache.GetPolicyResponse(r.Context(), req.Digest); err == nil && cachedResp != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedResp)
		return
	}

	// 2. Gather context from DB
	ctxData, err := h.db.GatherContext(r.Context(), req.Digest)
	if err != nil {
		// Log error and return standard response. In a real app we'd handle not found vs internal error
		http.Error(w, "failed to gather context for digest", http.StatusInternalServerError)
		return
	}

	// Convert context to map for OPA
	ctxMap := map[string]interface{}{}
	bytes, _ := json.Marshal(ctxData)
	json.Unmarshal(bytes, &ctxMap)

	// 3. Evaluate Policy
	allowed, reason, err := h.opa.Evaluate(r.Context(), ctxMap)
	if err != nil {
		http.Error(w, "policy evaluation failed", http.StatusInternalServerError)
		return
	}

	resp := domain.PolicyResponse{
		Digest:  req.Digest,
		Allowed: allowed,
		Reason:  reason,
		Cached:  false,
	}

	// 4. Update Cache
	_ = h.cache.SetPolicyResponse(r.Context(), req.Digest, resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
