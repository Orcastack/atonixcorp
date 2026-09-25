package metrics

import (
	"syscall"
)

type DiskCollector struct {
	Path string
}

func NewDiskCollector(path string) *DiskCollector {
	return &DiskCollector{Path: path}
}

func (d *DiskCollector) Collect() ([]Metric, error) {
	var stat syscall.Statfs_t
	_ = syscall.Statfs(d.Path, &stat)

	total := float64(stat.Blocks) * float64(stat.Bsize)
	free := float64(stat.Bfree) * float64(stat.Bsize)
	used := total - free

	metrics := []Metric{
		{
			Name:   "disk.total_bytes",
			Value:  total,
			Labels: map[string]string{"path": d.Path},
		},
		{
			Name:   "disk.used_bytes",
			Value:  used,
			Labels: map[string]string{"path": d.Path},
		},
		{
			Name:   "disk.free_bytes",
			Value:  free,
			Labels: map[string]string{"path": d.Path},
		},
	}

	return metrics, nil
}
