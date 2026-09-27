package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *APIServer) GetMetrics(w http.ResponseWriter, r *http.Request) {
	instanceID := chi.URLParam(r, "instance_id")

	metrics, err := s.metrics.CollectInstanceMetrics(instanceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(metrics)
}
