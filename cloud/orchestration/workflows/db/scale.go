package db

import (
	"fmt"
	"time"

	"atonixcorp/cloud/orchestration/events"
	"atonixcorp/cloud/orchestration/state"
)

func Scale(id string, cpu int, ram int) error {
	st := state.NewStateManager()
	ev := events.NewRouter()

	// Desired state
	st.SetDesired(id, "scaled")

	// Backend call (placeholder)
	fmt.Println("Scaling DB", id, "to CPU:", cpu, "RAM:", ram)

	// Actual state
	st.SetActual(id, "scaled")

	// Emit event
	ev.Emit(events.Event{
		Type:      "db.scale",
		Message:   fmt.Sprintf("DB %s scaled to %d CPU / %d MB RAM", id, cpu, ram),
		Timestamp: time.Now().Unix(),
	})

	return nil
}
