package compute

import (
	"fmt"
	"time"

	"atonixcorp/cloud/orchestration/events"
	"atonixcorp/cloud/orchestration/scheduler"
	"atonixcorp/cloud/orchestration/state"
)

func Migrate(id string) error {
	sched := scheduler.NewScheduler()
	st := state.NewStateManager()
	ev := events.NewRouter()

	// 1. Pick new node
	newNode, err := sched.PickNode(2, 4096)
	if err != nil {
		return err
	}

	// 2. Hypervisor call (placeholder)
	fmt.Println("Migrating instance", id, "to node", newNode)

	// 3. Actual state update
	st.SetActualNode(id, newNode)

	// 4. Emit event
	ev.Emit(events.Event{
		Type:      "compute.migrate",
		Message:   fmt.Sprintf("Instance %s migrated to %s", id, newNode),
		Timestamp: time.Now().Unix(),
	})

	return nil
}
