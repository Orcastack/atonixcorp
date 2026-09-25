package models

type WorkloadType string

const (
	WorkloadBatch   WorkloadType = "batch"
	WorkloadService WorkloadType = "service"
)

type Workload struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Type      WorkloadType `json:"type"`
	ProjectID string       `json:"project_id"`
	ImageID   string       `json:"image_id"`
	FlavorID  string       `json:"flavor_id"`
	NetworkID string       `json:"network_id"`
}
