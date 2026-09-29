package storage

import (
	"fmt"
	"time"

	"atonixcorp/cloud/orchestration/events"
	"atonixcorp/cloud/orchestration/state"
)

func Detach(volumeID string, instanceID string) error {
	st := state.NewStateManager()
	ev := events.NewRouter()

	// 1. Desired state
	st.SetDesired(volumeID, "detached")

	// 2. Backend call (placeholder)
	fmt.Println("Detaching volume", volumeID, "from instance", instanceID)

	// 3. Actual state
	st.SetActual(volumeID, "detached")

	// 4. Emit event
	ev.Emit(events.Event{
		Type:      "storage.detach",
		Message:   fmt.Sprintf("Volume %s detached from instance %s", volumeID, instanceID),
		Timestamp: time.Now().Unix(),
	})

	return nil
}
