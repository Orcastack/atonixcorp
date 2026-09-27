package services

import (
	"fmt"

	"atonixcorp/cloud/database/core"
	"atonixcorp/cloud/database/internal"
)

type ProvisioningService struct {
	controller *core.Controller
}

func NewProvisioningService(controller *core.Controller) *ProvisioningService {
	return &ProvisioningService{controller}
}

func (p *ProvisioningService) GetInstance(instanceID string) (*core.DBInstance, error) {
	return p.controller.GetInstance(instanceID)
}

// ProvisionDBInstance is the high-level workflow used by API/UI.
func (p *ProvisioningService) ProvisionDBInstance(req core.CreateInstanceRequest) (*core.DBInstance, error) {
	internal.Info("Provisioning DB instance for tenant %s", req.TenantID)

	inst, err := p.controller.CreateInstance(req)
	if err != nil {
		internal.Error("Provisioning failed: %v", err)
		return nil, fmt.Errorf("failed to provision instance: %w", err)
	}

	internal.Info("DB instance %s provisioned successfully", inst.ID)
	return inst, nil
}
