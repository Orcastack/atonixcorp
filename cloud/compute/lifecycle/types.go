package lifecycle

type VMAction string

const (
	ActionStart   VMAction = "start"
	ActionStop    VMAction = "stop"
	ActionReboot  VMAction = "reboot"
	ActionPause   VMAction = "pause"
	ActionMigrate VMAction = "migrate"
)

type VM struct {
	ID           string
	NodeID       string
	TargetNodeID string // for migrate
}
