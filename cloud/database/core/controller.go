package core

import (
	"fmt"
	"time"

	"atonixcorp/cloud/database/internal"
)

// Store interface (your Postgres metadata DB)
type Store interface {
	CreateInstance(*DBInstance) error
	UpdateInstance(*DBInstance) error
	GetInstance(id string) (*DBInstance, error)
}

// OpenStackClient interface (your openstack/ package)
type OpenStackClient interface {
	CreateServer(name string, cpu int, memoryMB int, networkID string, sgID string) (string, error)
	CreateVolume(name string, sizeGB int) (string, error)
	AttachVolume(serverID, volumeID string) error
	GetServerIP(serverID string) (string, error)
}

type Controller struct {
	store  Store
	os     OpenStackClient
	policy *Policy
	plans  map[string]Plan
}

func NewController(store Store, os OpenStackClient) *Controller {
	return &Controller{
		store:  store,
		os:     os,
		policy: NewPolicy(),
	}
}

func (c *Controller) GetInstance(id string) (*DBInstance, error) {
	return c.store.GetInstance(id)
}

func (c *Controller) SetPlans(plans map[string]Plan) {
	c.plans = plans
}

func (c *Controller) CreateInstance(req CreateInstanceRequest) (*DBInstance, error) {
	// 1. Validate request
	if err := c.policy.ValidateCreate(req); err != nil {
		return nil, err
	}

	// 2. Create metadata object
	inst := &DBInstance{
		ID:        internal.GenerateID("db"),
		TenantID:  req.TenantID,
		Name:      req.Name,
		Engine:    req.Engine,
		Plan:      req.Plan,
		CPU:       req.CPU,
		MemoryMB:  req.MemoryMB,
		StorageGB: req.StorageGB,
		Region:    req.Region,
		Status:    "provisioning",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 3. Save initial metadata
	if err := c.store.CreateInstance(inst); err != nil {
		return nil, fmt.Errorf("failed to save instance metadata: %w", err)
	}

	// 4. Provision OpenStack server
	serverID, err := c.os.CreateServer(inst.Name, inst.CPU, inst.MemoryMB, "default-network", "default-sg")
	if err != nil {
		inst.Status = "error"
		c.store.UpdateInstance(inst)
		return nil, fmt.Errorf("openstack server creation failed: %w", err)
	}
	inst.OpenStackServerID = serverID

	// 5. Provision volume
	volumeID, err := c.os.CreateVolume(inst.Name+"-data", inst.StorageGB)
	if err != nil {
		inst.Status = "error"
		c.store.UpdateInstance(inst)
		return nil, fmt.Errorf("openstack volume creation failed: %w", err)
	}
	inst.OpenStackVolumeID = volumeID

	// 6. Attach volume
	if err := c.os.AttachVolume(serverID, volumeID); err != nil {
		inst.Status = "error"
		c.store.UpdateInstance(inst)
		return nil, fmt.Errorf("volume attach failed: %w", err)
	}

	// 7. Get server IP
	ip, err := c.os.GetServerIP(serverID)
	if err != nil {
		inst.Status = "error"
		c.store.UpdateInstance(inst)
		return nil, fmt.Errorf("failed to get server IP: %w", err)
	}
	inst.Endpoint = fmt.Sprintf("%s:5432", ip)

	// 8. Mark instance running
	inst.Status = "running"
	inst.UpdatedAt = time.Now()
	if err := c.store.UpdateInstance(inst); err != nil {
		return nil, fmt.Errorf("failed to update instance metadata: %w", err)
	}

	return inst, nil
}
