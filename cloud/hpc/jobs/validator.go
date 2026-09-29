package jobs

import "errors"

var (
	ErrInvalidJobSpec = errors.New("invalid job specification")
)

type CreateJobRequest struct {
	TenantID   string
	UserID     string
	Name       string
	Backend    string
	Partition  string
	Account    string
	ScriptPath string
	CPUs       int
	MemoryMB   int
	GPUs       int
	WorkingDir string
	Env        map[string]string
}

func ValidateCreateJob(req CreateJobRequest) error {
	if req.TenantID == "" || req.UserID == "" {
		return ErrInvalidJobSpec
	}
	if req.Name == "" {
		return ErrInvalidJobSpec
	}
	if req.ScriptPath == "" {
		return ErrInvalidJobSpec
	}
	if req.CPUs <= 0 {
		return ErrInvalidJobSpec
	}
	if req.MemoryMB <= 0 {
		return ErrInvalidJobSpec
	}
	if req.Backend == "" {
		return ErrInvalidJobSpec
	}
	return nil
}
