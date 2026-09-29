package network

import (
	"fmt"
	"time"

	"atonixcorp/cloud/orchestration/events"
	"atonixcorp/cloud/orchestration/state"
)

func Create(id string, name string, cidr string) error {
	st := state.NewStateManager()
	ev := events.NewRouter()

	// 1. Desired state
	st.SetDesired(id, "created")

	// 2. Backend call (placeholder)
	fmt.Println("Creating network:", name, "CIDR:", cidr)

	// 3. Actual state
	st.SetActual(id, "created")

	// 4. Emit event
	ev.Emit(events.Event{
		Type:      "network.create",
		Message:   fmt.Sprintf("Network %s (%s) created", name, cidr),
		Timestamp: time.Now().Unix(),
	})

	return nil
}
