package orchestrator

import (
    "fmt"

    "atonixcorp/cloud/controlplane/bus"
    "atonixcorp/cloud/controlplane/models"
    "atonixcorp/cloud/compute/lifecycle"
    "atonixcorp/cloud/compute/rpc"
)

type Orchestrator struct {
    bus       *bus.Bus
    lifecycle *lifecycle.Manager
}

func NewOrchestrator(b *bus.Bus, rpcBackend rpc.RPC) *Orchestrator {
    return &Orchestrator{
        bus:       b,
        lifecycle: lifecycle.NewManager(rpcBackend),
    }
}

// ------------------------------------------------------------
// 1. Submit workload (entry point from API/CLI)
// ------------------------------------------------------------
func (o *Orchestrator) SubmitWorkload(wl models.Workload) {
    fmt.Println("Orchestrator: workload submitted:", wl.ID)

    o.bus.Publish(bus.Event{
        Type: "workload.submitted",
        Data: wl,
    })
}

// ------------------------------------------------------------
// 2. Start VM after scheduler picks a node
// ------------------------------------------------------------
func (o *Orchestrator) StartVM(vmID, nodeID string) {
    vm := lifecycle.VM{
        ID:     vmID,
        NodeID: nodeID,
    }

    fmt.Println("Orchestrator: start VM", vmID, "on node", nodeID)
    _ = o.lifecycle.Start(vm)
}

// ------------------------------------------------------------
// 3. Stop VM
// ------------------------------------------------------------
func (o *Orchestrator) StopVM(vmID, nodeID string) {
    vm := lifecycle.VM{
        ID:     vmID,
        NodeID: nodeID,
    }

    fmt.Println("Orchestrator: stop VM", vmID)
    _ = o.lifecycle.Stop(vm)
}

// ------------------------------------------------------------
// 4. Reboot VM
// ------------------------------------------------------------
func (o *Orchestrator) RebootVM(vmID, nodeID string) {
    vm := lifecycle.VM{
        ID:     vmID,
        NodeID: nodeID,
    }

    fmt.Println("Orchestrator: reboot VM", vmID)
    _ = o.lifecycle.Reboot(vm)
}

// ------------------------------------------------------------
// 5. Pause VM
// ------------------------------------------------------------
func (o *Orchestrator) PauseVM(vmID, nodeID string) {
    vm := lifecycle.VM{
        ID:     vmID,
        NodeID: nodeID,
    }

    fmt.Println("Orchestrator: pause VM", vmID)
    _ = o.lifecycle.Pause(vm)
}

// ------------------------------------------------------------
// 6. Migrate VM (node → node)
// ------------------------------------------------------------
func
