package handlers

import (
	"encoding/json"
	"net/http"

	"atonixcorp/cloud/hpc/nodes"
)

type NodeListHandler struct {
	Inventory *nodes.InventoryService
}

func NewNodeListHandler(inv *nodes.InventoryService) *NodeListHandler {
	return &NodeListHandler{Inventory: inv}
}

func (h *NodeListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ns, err := h.Inventory.ListAllNodes(r.Context())
	if err != nil {
		http.Error(w, "failed to list nodes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ns)
}
