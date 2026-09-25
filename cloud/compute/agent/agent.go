package agent

import (
	"atonixcorp/cloud/compute/rpc"
	"fmt"
)

type Agent struct {
	NodeID string
	RPC    rpc.RPC
}

func NewAgent(nodeID string, backend rpc.RPC) *Agent {
	return &Agent{
		NodeID: nodeID,
		RPC:    backend,
	}
}

func (a *Agent) Start() {
	fmt.Println("Compute Agent started on node:", a.NodeID)

	for {
		msg, err := a.RPC.Receive()
		if err != nil {
			fmt.Println("Agent receive error:", err)
			continue
		}

		a.HandleMessage(msg)
	}
}
