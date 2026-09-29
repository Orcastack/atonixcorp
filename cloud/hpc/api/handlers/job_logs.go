package handlers

import (
	"encoding/json"
	"net/http"

	"atonixcorp/cloud/hpc/jobs"

	"github.com/go-chi/chi/v5"
)

type JobLogsHandler struct {
	Logs *jobs.LogService
}

func NewJobLogsHandler(l *jobs.LogService) *JobLogsHandler {
	return &JobLogsHandler{Logs: l}
}

func (h *JobLogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing job id", http.StatusBadRequest)
		return
	}

	logs, err := h.Logs.List(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
