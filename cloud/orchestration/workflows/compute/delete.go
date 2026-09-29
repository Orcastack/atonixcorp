package compute

import (
	"fmt"
	"time"

	"atonixcorp/cloud/orchestration/events"
	"atonixcorp/cloud/orchestration/state"
)

func Delete(id string) error {
	st := state.NewStateManager()
	ev := events.NewRouter()

	// 1. Desired state
	st.SetDesired(id, "deleted")

	// 2. Hypervisor call (placeholder)
	fmt.Println("Deleting instance", id)

	// 3. Actual state
	st.SetActual(id, "deleted")

	// 4. Emit event
	ev.Emit(events.Event{
		Type:      "compute.delete",
		Message:   fmt.Sprintf("Instance %s deleted", id),
		Timestamp: time.Now().Unix(),
	})

	return nil
}
