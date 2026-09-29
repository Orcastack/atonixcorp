package monitoring

import "fmt"

// DatasetMetrics tracks dataset usage, access frequency, and size.
type DatasetMetrics struct {
	DatasetID      string
	AccessCount    int
	LastAccessTime string
	SizeBytes      int64
}

func NewDatasetMetrics(datasetID string) *DatasetMetrics {
	return &DatasetMetrics{
		DatasetID: datasetID,
	}
}

func (m *DatasetMetrics) RecordAccess(timestamp string) {
	m.AccessCount++
	m.LastAccessTime = timestamp
	fmt.Printf("Dataset accessed: %s at %s\n", m.DatasetID, timestamp)
}

func (m *DatasetMetrics) UpdateSize(bytes int64) {
	m.SizeBytes = bytes
	fmt.Printf("Dataset size updated: %s = %d bytes\n", m.DatasetID, bytes)
}

func (m *DatasetMetrics) Snapshot() map[string]any {
	return map[string]any{
		"dataset":      m.DatasetID,
		"access_count": m.AccessCount,
		"last_access":  m.LastAccessTime,
		"size_bytes":   m.SizeBytes,
	}
}
