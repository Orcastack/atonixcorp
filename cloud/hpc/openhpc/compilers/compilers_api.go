package compilers

// API is the cloud-facing compiler service API.
type API struct {
	svc *Service
}

func NewAPI() *API {
	return &API{
		svc: NewService(),
	}
}

// ListToolchains returns all compiler toolchains available.
func (api *API) ListToolchains() []Compiler {
	return api.svc.List()
}

// ConfigureBuild prepares compiler runtime for a user build/job.
func (api *API) ConfigureBuild(compilerName, language, opt string, targetGPU bool) (map[string]string, error) {
	cfg := RuntimeConfig{
		CompilerName: compilerName,
		Language:     language,
		Optimization: opt,
		TargetGPU:    targetGPU,
	}
	return BuildRuntime(cfg, api.svc)
}
