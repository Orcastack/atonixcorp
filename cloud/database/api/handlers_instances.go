package api

import (
	"encoding/json"
	"net/http"

	"atonixcorp/cloud/database/core"

	"github.com/go-chi/chi/v5"
)

func (s *APIServer) CreateInstance(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenant_id")

	var req core.CreateInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	req.TenantID = tenantID

	inst, err := s.provisioning.ProvisionDBInstance(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(inst)
}

func (s *APIServer) GetInstance(w http.ResponseWriter, r *http.Request) {
	instanceID := chi.URLParam(r, "instance_id")

	inst, err := s.provisioning.GetInstance(instanceID)
	if err != nil {
		http.Error(w, "instance not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(inst)
}
