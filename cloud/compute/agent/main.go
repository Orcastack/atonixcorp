package main

import (
	"atonixcorp/cloud/compute/agent"
	"atonixcorp/cloud/compute/rpc"
)

func main() {
	// Choose backend (local, nats, rabbitmq, grpc, zeromq)
	backend := rpc.NewLocalRPC()

	a := agent.NewAgent("node-001", backend)

	a.Register()
	a.ReportHealth()
	a.Start()
}
