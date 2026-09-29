package network

import (
	"fmt"
	"time"

	"atonixcorp/cloud/orchestration/events"
	"atonixcorp/cloud/orchestration/state"
)

func AddRoute(id string, destination string, nextHop string) error {
	st := state.NewStateManager()
	ev := events.NewRouter()

	// 1. Desired state
	st.SetDesired(id, "route-updated")

	// 2. Backend call (placeholder)
	fmt.Println("Adding route:", destination, "via", nextHop)

	// 3. Actual state
	st.SetActual(id, "route-updated")

	// 4. Emit event
	ev.Emit(events.Event{
		Type:      "network.route",
		Message:   fmt.Sprintf("Route added: %s → %s", destination, nextHop),
		Timestamp: time.Now().Unix(),
	})

	return nil
}
