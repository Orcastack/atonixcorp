package compute

import (
	"fmt"
	"time"

	"atonixcorp/cloud/orchestration/events"
	"atonixcorp/cloud/orchestration/state"
)

func Stop(id string) error {
	st := state.NewStateManager()
	ev := events.NewRouter()

	// 1. Desired state
	st.SetDesired(id, "stopped")

	// 2. Hypervisor call (placeholder)
	fmt.Println("Stopping instance", id)

	// 3. Actual state
	st.SetActual(id, "stopped")

	// 4. Emit event
	ev.Emit(events.Event{
		Type:      "compute.stop",
		Message:   fmt.Sprintf("Instance %s stopped", id),
		Timestamp: time.Now().Unix(),
	})

	return nil
}
