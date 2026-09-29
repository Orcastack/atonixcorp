package handlers

import (
	"encoding/json"
	"net/http"

	"atonixcorp/cloud/hpc/nodes"
)

type PartitionListHandler struct {
	Inventory *nodes.InventoryService
}

func NewPartitionListHandler(inv *nodes.InventoryService) *PartitionListHandler {
	return &PartitionListHandler{Inventory: inv}
}

func (h *PartitionListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps, err := h.Inventory.ListPartitions(r.Context())
	if err != nil {
		http.Error(w, "failed to list partitions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ps)
}
