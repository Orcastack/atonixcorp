package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type BackupRequest struct {
	Type string `json:"type"` // full | snapshot
}

func (s *APIServer) TriggerBackup(w http.ResponseWriter, r *http.Request) {
	instanceID := chi.URLParam(r, "instance_id")

	var req BackupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	backup, err := s.backups.TriggerBackup(instanceID, req.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(backup)
}
