package monitoring

import "fmt"

// HealthStatus tracks dataset and storage backend health.
type HealthStatus struct {
	DatasetID string
	StorageOK bool
	IOK       bool
	LastCheck string
	Message   string
}

func NewHealthStatus(datasetID string) *HealthStatus {
	return &HealthStatus{
		DatasetID: datasetID,
	}
}

func (h *HealthStatus) Update(storageOK, ioOK bool, timestamp, msg string) {
	h.StorageOK = storageOK
	h.IOK = ioOK
	h.LastCheck = timestamp
	h.Message = msg

	fmt.Printf("Health updated for %s: storage=%v io=%v msg=%s\n",
		h.DatasetID, storageOK, ioOK, msg)
}

func (h *HealthStatus) Snapshot() map[string]any {
	return map[string]any{
		"dataset":    h.DatasetID,
		"storage_ok": h.StorageOK,
		"io_ok":      h.IOK,
		"last_check": h.LastCheck,
		"message":    h.Message,
	}
}
