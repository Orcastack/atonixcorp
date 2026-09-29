package storage

import (
	"fmt"
	"time"

	"atonixcorp/cloud/orchestration/events"
	"atonixcorp/cloud/orchestration/state"
)

func Attach(volumeID string, instanceID string) error {
	st := state.NewStateManager()
	ev := events.NewRouter()

	// 1. Desired state
	st.SetDesired(volumeID, "attached")

	// 2. Backend call (placeholder)
	fmt.Println("Attaching volume", volumeID, "to instance", instanceID)

	// 3. Actual state
	st.SetActual(volumeID, "attached")

	// 4. Emit event
	ev.Emit(events.Event{
		Type:      "storage.attach",
		Message:   fmt.Sprintf("Volume %s attached to instance %s", volumeID, instanceID),
		Timestamp: time.Now().Unix(),
	})

	return nil
}
