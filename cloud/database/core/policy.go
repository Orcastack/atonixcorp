package core

import "fmt"

type Policy struct{}

func NewPolicy() *Policy {
	return &Policy{}
}

func (p *Policy) ValidateCreate(req CreateInstanceRequest) error {
	if req.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}

	switch req.Engine {
	case EnginePostgres, EngineTimescale, EngineMongo, EngineRedis:
	default:
		return fmt.Errorf("unsupported engine: %s", req.Engine)
	}

	if req.CPU < 1 {
		return fmt.Errorf("cpu must be >= 1")
	}
	if req.MemoryMB < 1024 {
		return fmt.Errorf("memory must be >= 1024 MB")
	}
	if req.StorageGB < 10 {
		return fmt.Errorf("storage must be >= 10 GB")
	}

	return nil
}
