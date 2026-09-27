package core

import "time"

// Supported engines
const (
	EnginePostgres  = "postgres"
	EngineTimescale = "timescale"
	EngineMongo     = "mongo"
	EngineRedis     = "redis"
)

type Plan struct {
	CPU          int         `yaml:"cpu"`
	MemoryMB     int         `yaml:"memory_mb"`
	StorageGB    int         `yaml:"storage_gb"`
	HA           bool        `yaml:"ha"`
	Backups      string      `yaml:"backups"`
	ReadReplicas int         `yaml:"read_replicas"`
	Engines      []string    `yaml:"engines"`
	PriceUSD     float64     `yaml:"price_usd"`
	NodePricing  NodePricing `yaml:"node_pricing"`
}

type NodePricing struct {
	CPUUSD     float64 `yaml:"cpu_usd"`
	MemoryUSD  float64 `yaml:"memory_usd"`
	StorageUSD float64 `yaml:"storage_usd"`
}

// DBInstance represents a database instance provisioned on OpenStack.
type DBInstance struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	Name      string `json:"name"`
	Engine    string `json:"engine"`
	Plan      string `json:"plan"`
	CPU       int    `json:"cpu"`
	MemoryMB  int    `json:"memory_mb"`
	StorageGB int    `json:"storage_gb"`
	Region    string `json:"region"`
	Status    string `json:"status"`
	Endpoint  string `json:"endpoint"`

	// OpenStack resources
	OpenStackServerID  string `json:"openstack_server_id"`
	OpenStackVolumeID  string `json:"openstack_volume_id"`
	OpenStackNetworkID string `json:"openstack_network_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Backup metadata
type Backup struct {
	ID         string    `json:"id"`
	InstanceID string    `json:"instance_id"`
	Type       string    `json:"type"`     // full | snapshot
	Location   string    `json:"location"` // s3://...
	CreatedAt  time.Time `json:"created_at"`
}

type DBMetrics struct {
	InstanceID   string    `json:"instance_id"`
	CPUUsage     float64   `json:"cpu_usage"`
	MemoryUsage  float64   `json:"memory_usage"`
	StorageUsage float64   `json:"storage_usage"`
	Connections  int       `json:"connections"`
	Timestamp    time.Time `json:"timestamp"`
}

// CreateInstanceRequest is the payload from API/UI.
type CreateInstanceRequest struct {
	TenantID  string `json:"tenant_id"`
	Name      string `json:"name"`
	Engine    string `json:"engine"`
	Plan      string `json:"plan"`
	CPU       int    `json:"cpu"`
	MemoryMB  int    `json:"memory_mb"`
	StorageGB int    `json:"storage_gb"`
	Region    string `json:"region"`
}
