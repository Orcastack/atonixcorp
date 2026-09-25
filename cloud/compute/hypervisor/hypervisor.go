package hypervisor

type VMState string

const (
	VMRunning VMState = "running"
	VMStopped VMState = "stopped"
	VMPaused  VMState = "paused"
)

type Hypervisor interface {
	StartVM(id string) error
	StopVM(id string) error
	RebootVM(id string) error
	GetState(id string) (VMState, error)
}
