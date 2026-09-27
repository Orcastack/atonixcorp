package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"context"

	sdk "atonixcorp/cloud/devtools/sdk/go"
)

func ListProjects(w http.ResponseWriter, r *http.Request) {
	client := sdk.NewClient(sdk.ClientOptions{
		Endpoint: os.Getenv("ATCLOUD_ENDPOINT"),
		Token:    os.Getenv("ATCLOUD_TOKEN"),
	})

	projects, err := client.ListProjects(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}
