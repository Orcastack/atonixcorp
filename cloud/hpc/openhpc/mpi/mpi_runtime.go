package mpi

import "fmt"

// MPIRuntimeConfig defines how MPI is configured for a job.
type MPIRuntimeConfig struct {
	Implementation string
	Nodes          []string
	TasksPerNode   int
}

// BuildRuntime generates environment variables for the scheduler.
func BuildRuntime(cfg MPIRuntimeConfig, svc *MPIService) (map[string]string, error) {
	impl := svc.Get(cfg.Implementation)
	if impl == nil {
		return nil, fmt.Errorf("MPI implementation '%s' not found", cfg.Implementation)
	}

	env := map[string]string{
		"MPI_IMPL":           impl.Name,
		"MPI_VERSION":        impl.Version,
		"MPI_MODULE":         impl.Module,
		"MPI_CONTAINER":      impl.Container,
		"MPI_NODES":          fmt.Sprintf("%v", cfg.Nodes),
		"MPI_TASKS_PER_NODE": fmt.Sprintf("%d", cfg.TasksPerNode),
	}

	return env, nil
}
