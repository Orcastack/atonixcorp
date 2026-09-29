package mpi

// MPIImplementation represents an MPI stack provided by the cloud.
type MPIImplementation struct {
	Name        string
	Version     string
	Description string
	Container   string // cloud container image
	Module      string // environment module
}

// MPIService exposes MPI runtimes to cloud users.
type MPIService struct {
	Implementations []MPIImplementation
}

// NewMPIService initializes the cloud MPI service.
func NewMPIService() *MPIService {
	return &MPIService{
		Implementations: []MPIImplementation{
			{
				Name:        "openmpi",
				Version:     "4.1",
				Description: "OpenMPI runtime for distributed HPC workloads",
				Container:   "registry.atonixcorp.local/mpi/openmpi:4.1",
				Module:      "mpi/openmpi/4.1",
			},
			{
				Name:        "mpich",
				Version:     "3.4",
				Description: "MPICH runtime for scientific workloads",
				Container:   "registry.atonixcorp.local/mpi/mpich:3.4",
				Module:      "mpi/mpich/3.4",
			},
		},
	}
}

// List returns all MPI runtimes available on the cloud.
func (s *MPIService) List() []MPIImplementation {
	return s.Implementations
}

// Get returns a specific MPI runtime.
func (s *MPIService) Get(name string) *MPIImplementation {
	for _, impl := range s.Implementations {
		if impl.Name == name {
			return &impl
		}
	}
	return nil
}
