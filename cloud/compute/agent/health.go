package agent

import (
	"atonixcorp/cloud/compute/rpc"
	"fmt"
	"runtime"
)

func (a *Agent) ReportHealth() {
	health := map[string]any{
		"node":        a.NodeID,
		"cpus":        runtime.NumCPU(),
		"go_routines": runtime.NumGoroutine(),
	}

	a.RPC.Send(rpc.RPCMessage{
		Type:    "node.health",
		Payload: health,
	})

	fmt.Println("Agent: health sent")
}
