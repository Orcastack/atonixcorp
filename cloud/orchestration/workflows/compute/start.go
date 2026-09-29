package compute

import (
	"fmt"
	"time"

	"atonixcorp/cloud/orchestration/events"
	"atonixcorp/cloud/orchestration/scheduler"
	"atonixcorp/cloud/orchestration/state"
)

func Start(id string) error {
	sched := scheduler.NewScheduler()
	st := state.NewStateManager()
	ev := events.NewRouter()

	// 1. Desired state
	st.SetDesired(id, "running")

	// 2. Pick node
	node, err := sched.PickNode(2, 4096)
	if err != nil {
		return err
	}

	// 3. Hypervisor call (placeholder)
	fmt.Println("Starting instance", id, "on node", node)

	// 4. Actual state
	st.SetActual(id, "running")

	// 5. Emit event
	ev.Emit(events.Event{
		Type:      "compute.start",
		Message:   fmt.Sprintf("Instance %s started on %s", id, node),
		Timestamp: time.Now().Unix(),
	})

	return nil
}
