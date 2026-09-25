package agent

import (
	"atonixcorp/cloud/compute/rpc"
	"fmt"
)

func (a *Agent) Register() {
	a.RPC.Send(rpc.RPCMessage{
		Type: "node.register",
		Payload: map[string]any{
			"node_id": a.NodeID,
		},
	})

	fmt.Println("Agent: registered with control plane")
}
