package handlers

import (
	"encoding/json"
	"net/http"
)

func Version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"portal":  "ATCloud DevTools Portal v1",
		"backend": "go",
	})
}
