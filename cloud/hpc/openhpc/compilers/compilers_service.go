package compilers

// Compiler represents a compiler toolchain provided by the cloud.
type Compiler struct {
	Name        string // e.g. "gcc", "clang", "nvhpc"
	Version     string
	Description string
	Container   string // container image
	Module      string // environment module
	Languages   []string
}

// Service exposes compilers as cloud services.
type Service struct {
	Compilers []Compiler
}

// NewService initializes the compiler service with all toolchains.
func NewService() *Service {
	return &Service{
		Compilers: []Compiler{
			{
				Name:        "gcc",
				Version:     "12",
				Description: "GNU Compiler Collection for C/C++/Fortran",
				Container:   "registry.atonixcorp.local/compilers/gcc:12",
				Module:      "compiler/gcc/12",
				Languages:   []string{"c", "c++", "fortran"},
			},
			{
				Name:        "clang",
				Version:     "16",
				Description: "Clang/LLVM toolchain",
				Container:   "registry.atonixcorp.local/compilers/clang:16",
				Module:      "compiler/clang/16",
				Languages:   []string{"c", "c++"},
			},
			{
				Name:        "nvhpc",
				Version:     "24.3",
				Description: "NVIDIA HPC SDK for GPU-accelerated codes",
				Container:   "registry.atonixcorp.local/compilers/nvhpc:24.3",
				Module:      "compiler/nvhpc/24.3",
				Languages:   []string{"c", "c++", "fortran"},
			},
		},
	}
}

// List returns all compiler toolchains.
func (s *Service) List() []Compiler {
	return s.Compilers
}

// Get returns a specific compiler by name.
func (s *Service) Get(name string) *Compiler {
	for _, c := range s.Compilers {
		if c.Name == name {
			return &c
		}
	}
	return nil
}
