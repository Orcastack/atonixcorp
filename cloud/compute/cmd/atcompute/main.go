package main

import (
	"atonixcorp/cloud/compute/agent"
	"atonixcorp/cloud/compute/rpc"
)

func main() {
	backend := rpc.NewLocalRPC()

	nodeAgent := agent.NewAgent("node-001", backend)
	nodeAgent.Register()
	nodeAgent.ReportHealth()
	nodeAgent.Start()
}