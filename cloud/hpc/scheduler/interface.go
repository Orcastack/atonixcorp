package scheduler

import (
	"context"
)

// JobSpec defines what the control plane expects for job submission.
type JobSpec struct {
	JobName    string
	ScriptPath string
	Partition  string
	Account    string
	TimeLimit  string // "01:00:00"
	CPUs       int
	MemoryMB   int
	GPUs       int
	WorkingDir string
	Env        map[string]string
}

// JobStatus represents normalized job status across Slurm/K8s.
type JobStatus struct {
	JobID     string
	Name      string
	State     string
	Partition string
	User      string
	Nodes     string
	Raw       string
}

// Scheduler is the unified interface for all HPC backends.
type Scheduler interface {
	Submit(ctx context.Context, spec JobSpec) (string, error)
	Status(ctx context.Context, jobID string) (*JobStatus, error)
	Cancel(ctx context.Context, jobID string) error
}
