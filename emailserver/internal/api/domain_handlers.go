package api

import (
	"encoding/json"
	"net/http"
)

func (d Dependencies) handleRegisterDomain(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := d.Domain.RegisterDomain(req.Name); err != nil {
		d.Logger.Error("register domain:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d.Analytics.SendEvent("domain_register", map[string]any{
		"user":   userID,
		"domain": req.Name,
	})

	w.WriteHeader(http.StatusCreated)
}

func (d Dependencies) handleListDomains(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	doms, err := d.Domain.ListDomains()
	if err != nil {
		d.Logger.Error("list domains:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doms)
}
