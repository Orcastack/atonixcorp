package scheduler

import (
	"fmt"

	"atonixcorp/cloud/controlplane/bus"
	"atonixcorp/cloud/controlplane/models"
	"atonixcorp/cloud/controlplane/orchestrator"
)

type Scheduler struct {
	bus          *bus.Bus
	orchestrator *orchestrator.Orchestrator
}

func NewScheduler(b *bus.Bus, orch *orchestrator.Orchestrator) *Scheduler {
	s := &Scheduler{
		bus:          b,
		orchestrator: orch,
	}

	b.Subscribe("workload.submitted", s.onWorkloadSubmitted)
	return s
}

// ------------------------------------------------------------
// 1. Handle workload submission
// ------------------------------------------------------------
func (s *Scheduler) onWorkloadSubmitted(e bus.Event) {
	wl, ok := e.Data.(models.Workload)
	if !ok {
		return
	}

	fmt.Println("Scheduler: workload received:", wl.ID)

	// ------------------------------------------------------------
	// 2. Pick a compute node (placeholder)
	// ------------------------------------------------------------
	nodeID := s.pickNode(wl)
	fmt.Println("Scheduler: selected node:", nodeID)

	// ------------------------------------------------------------
	// 3. Ask orchestrator to start VM on that node
	// ------------------------------------------------------------
	s.orchestrator.StartVM(wl.ID, nodeID)

	// ------------------------------------------------------------
	// 4. Publish scheduled event
	// ------------------------------------------------------------
	s.bus.Publish(bus.Event{
		Type: "workload.scheduled",
		Data: map[string]any{
			"vm_id":   wl.ID,
			"node_id": nodeID,
		},
	})
}

// ------------------------------------------------------------
// Node selection (placeholder)
// ------------------------------------------------------------
func (s *Scheduler) pickNode(wl models.Workload) string {
	// TODO: integrate real placement engine
	return "node-001"
}
