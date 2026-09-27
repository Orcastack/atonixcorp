package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type DevToolsAPI struct {
	Version     string            `json:"version"`
	Build       string            `json:"build"`
	Timestamp   int64             `json:"timestamp"`
	Services    map[string]string `json:"services"`
	SDKs        []string          `json:"sdks"`
	CLICommands []string          `json:"cli_commands"`
	Portal      string            `json:"portal"`
}

func Info(w http.ResponseWriter, r *http.Request) {
	resp := DevToolsAPI{
		Version:   "v1.0.0",
		Build:     "stable",
		Timestamp: time.Now().Unix(),

		Services: map[string]string{
			"compute":       "/compute",
			"network":       "/network",
			"storage":       "/storage",
			"swift":         "/swift",
			"identity":      "/identity",
			"images":        "/images",
			"orchestration": "/orchestration",
			"projects":      "/projects",
			"telemetry":     "/telemetry",
		},

		SDKs: []string{
			"go",
			"python",
			"javascript",
			"java",
		},

		CLICommands: []string{
			"create compute",
			"delete compute",
			"list compute",
			"create network",
			"delete network",
			"list network",
			"create storage",
			"delete storage",
			"list storage",
			"create container",
			"delete container",
			"upload object",
			"delete object",
			"identity",
			"list images",
			"list events",
			"list metrics",
			"list projects",
			"telemetry",
		},

		Portal: "/portal",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
