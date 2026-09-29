package monitoring

import "fmt"

// IOMetrics tracks parallel I/O performance for datasets.
type IOMetrics struct {
	DatasetID           string
	ReadThroughputMBps  float64
	WriteThroughputMBps float64
	LatencyMs           float64
	ParallelNodes       int
}

func NewIOMetrics(datasetID string) *IOMetrics {
	return &IOMetrics{
		DatasetID: datasetID,
	}
}

func (m *IOMetrics) Update(readMBps, writeMBps, latency float64, nodes int) {
	m.ReadThroughputMBps = readMBps
	m.WriteThroughputMBps = writeMBps
	m.LatencyMs = latency
	m.ParallelNodes = nodes

	fmt.Printf("I/O metrics updated for %s: read=%.2fMB/s write=%.2fMB/s latency=%.2fms nodes=%d\n",
		m.DatasetID, readMBps, writeMBps, latency, nodes)
}

func (m *IOMetrics) Snapshot() map[string]any {
	return map[string]any{
		"dataset":        m.DatasetID,
		"read_mb_s":      m.ReadThroughputMBps,
		"write_mb_s":     m.WriteThroughputMBps,
		"latency_ms":     m.LatencyMs,
		"parallel_nodes": m.ParallelNodes,
	}
}
