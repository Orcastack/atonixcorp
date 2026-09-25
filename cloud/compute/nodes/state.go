package nodes

type NodeState string

const (
	NodeOnline   NodeState = "online"
	NodeOffline  NodeState = "offline"
	NodeDraining NodeState = "draining"
)
