package state

type NodeState struct {
	ID       string `json:"id"`
	CPUFree  int    `json:"cpu_free"`
	RAMFree  int    `json:"ram_free"`
	DiskFree int    `json:"disk_free"`
}
