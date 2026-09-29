package mpi

import "fmt"

// JobMPIConfig represents MPI settings for a scheduled job.
type JobMPIConfig struct {
	ImplementationName string
	Nodes              []string
	TasksPerNode       int
}

// GenerateEnv generates environment hints for the scheduler / job launcher.
func GenerateEnv(cfg JobMPIConfig, svc *Service) (map[string]string, error) {
	impl := svc.Get(cfg.ImplementationName)
	if impl == nil {
		return nil, fmt.Errorf("MPI implementation '%s' not found", cfg.ImplementationName)
	}

	env := map[string]string{
		"MPI_IMPL":           impl.Name,
		"MPI_VERSION":        impl.Version,
		"MPI_MODULE":         impl.ModuleName,
		"MPI_CONTAINER":      impl.ContainerImage,
		"MPI_NODES":          fmt.Sprintf("%v", cfg.Nodes),
		"MPI_TASKS_PER_NODE": fmt.Sprintf("%d", cfg.TasksPerNode),
	}

	return env, nil
}
