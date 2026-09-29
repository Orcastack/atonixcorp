package network

import (
	"fmt"
	"time"

	"atonixcorp/cloud/orchestration/events"
	"atonixcorp/cloud/orchestration/state"
)

func ApplySecurity(id string, rule string, action string) error {
	st := state.NewStateManager()
	ev := events.NewRouter()

	// 1. Desired state
	st.SetDesired(id, "security-updated")

	// 2. Backend call (placeholder)
	fmt.Println("Applying security rule:", rule, "action:", action)

	// 3. Actual state
	st.SetActual(id, "security-updated")

	// 4. Emit event
	ev.Emit(events.Event{
		Type:      "network.security",
		Message:   fmt.Sprintf("Security rule applied: %s (%s)", rule, action),
		Timestamp: time.Now().Unix(),
	})

	return nil
}
