package agent

import (
	"encoding/json"
	"net/http"
)

type API struct {
	publisher *Publisher
}

func NewAPI(publisher *Publisher) *API {
	return &API{publisher: publisher}
}

func (a *API) HandleSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	var obs WorkloadObservation
	if err := json.NewDecoder(r.Body).Decode(&obs); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if err := a.publisher.Publish(r.Context(), obs); err != nil {
		http.Error(w, "failed to publish observation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "observation injected"})
}
