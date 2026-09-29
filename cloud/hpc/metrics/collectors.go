package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type HPCCollectors struct {
	JobSubmissions   *prometheus.CounterVec
	JobFailures      *prometheus.CounterVec
	JobDuration      *prometheus.HistogramVec
	ClusterCPUUsage  *prometheus.GaugeVec
	ClusterGPUUsage  *prometheus.GaugeVec
	ClusterNodeCount *prometheus.GaugeVec
}

func NewHPCCollectors() *HPCCollectors {
	return &HPCCollectors{
		JobSubmissions: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "atonix_hpc_job_submissions_total",
				Help: "Total number of HPC jobs submitted",
			},
			[]string{"tenant", "backend", "partition"},
		),
		JobFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "atonix_hpc_job_failures_total",
				Help: "Total number of failed HPC jobs",
			},
			[]string{"tenant", "backend", "partition"},
		),
		JobDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "atonix_hpc_job_duration_seconds",
				Help:    "Duration of completed HPC jobs",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"tenant", "backend", "partition"},
		),
		ClusterCPUUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "atonix_hpc_cluster_cpu_usage_ratio",
				Help: "Cluster CPU usage ratio (0–1)",
			},
			[]string{"cluster"},
		),
		ClusterGPUUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "atonix_hpc_cluster_gpu_usage_ratio",
				Help: "Cluster GPU usage ratio (0–1)",
			},
			[]string{"cluster"},
		),
		ClusterNodeCount: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "atonix_hpc_cluster_node_count",
				Help: "Number of nodes in the HPC cluster",
			},
			[]string{"cluster"},
		),
	}
}

func (c *HPCCollectors) Register(reg prometheus.Registerer) {
	reg.MustRegister(
		c.JobSubmissions,
		c.JobFailures,
		c.JobDuration,
		c.ClusterCPUUsage,
		c.ClusterGPUUsage,
		c.ClusterNodeCount,
	)
}
