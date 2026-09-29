package nodes

import (
	"time"
)

type NodeState string

const (
	NodeStateUnknown  NodeState = "unknown"
	NodeStateOnline   NodeState = "online"
	NodeStateOffline  NodeState = "offline"
	NodeStateDraining NodeState = "draining"
	NodeStateError    NodeState = "error"
)

type NodeRole string

const (
	NodeRoleCompute NodeRole = "compute"
	NodeRoleLogin   NodeRole = "login"
	NodeRoleStorage NodeRole = "storage"
	NodeRoleGPU     NodeRole = "gpu"
)

type Node struct {
	ID          string            `json:"id"`
	Hostname    string            `json:"hostname"`
	Address     string            `json:"address"`
	Role        NodeRole          `json:"role"`
	State       NodeState         `json:"state"`
	CPUs        int               `json:"cpus"`
	MemoryGB    int               `json:"memory_gb"`
	GPUs        int               `json:"gpus"`
	Labels      map[string]string `json:"labels"`
	Partitions  []string          `json:"partitions"`
	LastSeenAt  time.Time         `json:"last_seen_at"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	HealthScore float64           `json:"health_score"`
}

type Partition struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Nodes       []string          `json:"nodes"`
	Labels      map[string]string `json:"labels"`
	MaxCPUs     int               `json:"max_cpus"`
	MaxGPUs     int               `json:"max_gpus"`
	Active      bool              `json:"active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
