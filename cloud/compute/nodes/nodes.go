package nodes

type Node struct {
	ID         string
	Hostname   string
	CPUs       int
	RAMMB      int
	DiskGB     int
	Hypervisor string
	State      NodeState
}
