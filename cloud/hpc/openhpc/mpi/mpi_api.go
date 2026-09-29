package mpi

// MPIAPI is the cloud-facing API for users.
type MPIAPI struct {
	svc *MPIService
}

func NewMPIAPI() *MPIAPI {
	return &MPIAPI{
		svc: NewMPIService(),
	}
}

// ListRuntimes returns all MPI runtimes available.
func (api *MPIAPI) ListRuntimes() []MPIImplementation {
	return api.svc.List()
}

// ConfigureJobMPI prepares MPI runtime for a job.
func (api *MPIAPI) ConfigureJobMPI(impl string, nodes []string, tasks int) (map[string]string, error) {
	cfg := MPIRuntimeConfig{
		Implementation: impl,
		Nodes:          nodes,
		TasksPerNode:   tasks,
	}
	return BuildRuntime(cfg, api.svc)
}
