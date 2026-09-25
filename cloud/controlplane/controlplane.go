package controlplane

import (
	"atonixcorp/cloud/controlplane/audit"
	"atonixcorp/cloud/controlplane/bus"
	"atonixcorp/cloud/controlplane/orchestrator"
	"atonixcorp/cloud/controlplane/policy"
	"atonixcorp/cloud/controlplane/scheduler"

	"atonixcorp/cloud/compute/rpc"
	"atonixcorp/cloud/core/keystone"
	"atonixcorp/cloud/core/nova"
)

type ControlPlane struct {
	Bus          *bus.Bus
	Orchestrator *orchestrator.Orchestrator
	RPC          rpc.RPC
}

type Config struct {
	RPCBackend string
}

func NewControlPlane(cfg Config, ks *keystone.Client, nv *nova.Client) *ControlPlane {
	b := bus.NewBus()

	// ------------------------------------------------------------
	// SELECT RPC BACKEND HERE
	// ------------------------------------------------------------
	var rpcBackend rpc.RPC

	switch cfg.RPCBackend {
	case "nats":
		rpcBackend, _ = rpc.NewNatsRPC("nats://localhost:4222", "atonix.compute")
	case "rabbitmq":
		rpcBackend, _ = rpc.NewRabbitRPC("amqp://guest:guest@localhost:5672/", "atonix.queue")
	case "grpc":
		rpcBackend = rpc.NewGRPCRPC()
	case "zeromq":
		rpcBackend = rpc.NewZeroMQRPC()
	default:
		rpcBackend = rpc.NewLocalRPC()
	}

	// ------------------------------------------------------------
	// INITIALIZE CONTROL PLANE COMPONENTS
	// ------------------------------------------------------------
	_ = policy.NewEngine(b)
	_ = audit.NewService(b)
	_ = scheduler.NewScheduler(b, ks, nv, rpcBackend)

	orch := orchestrator.NewOrchestrator(b, rpcBackend)

	return &ControlPlane{
		Bus:          b,
		Orchestrator: orch,
		RPC:          rpcBackend,
	}
}
