package metrics

import (
	"runtime"
)

type CPUCollector struct{}

func NewCPUCollector() *CPUCollector {
	return &CPUCollector{}
}

func (c *CPUCollector) Collect() ([]Metric, error) {
	metrics := []Metric{
		{
			Name:  "cpu.cores",
			Value: float64(runtime.NumCPU()),
			Labels: map[string]string{
				"source": "node",
			},
		},
	}

	return metrics, nil
}
