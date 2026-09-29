package handlers

import (
	"encoding/json"
	"net/http"

	"atonixcorp/cloud/hpc/jobs"
)

type JobSubmitHandler struct {
	Service *jobs.Service
}

func NewJobSubmitHandler(s *jobs.Service) *JobSubmitHandler {
	return &JobSubmitHandler{Service: s}
}

func (h *JobSubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req jobs.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	job, err := h.Service.CreateAndSubmit(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}
