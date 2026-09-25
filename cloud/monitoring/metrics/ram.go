package metrics

import (
	"syscall"
)

type RAMCollector struct{}

func NewRAMCollector() *RAMCollector {
	return &RAMCollector{}
}

func (r *RAMCollector) Collect() ([]Metric, error) {
	var info syscall.Sysinfo_t
	_ = syscall.Sysinfo(&info)

	total := float64(info.Totalram)
	free := float64(info.Freeram)
	used := total - free

	metrics := []Metric{
		{
			Name:   "ram.total_bytes",
			Value:  total,
			Labels: map[string]string{"source": "node"},
		},
		{
			Name:   "ram.used_bytes",
			Value:  used,
			Labels: map[string]string{"source": "node"},
		},
		{
			Name:   "ram.free_bytes",
			Value:  free,
			Labels: map[string]string{"source": "node"},
		},
	}

	return metrics, nil
}
