package db

import (
	"fmt"
	"time"

	"atonixcorp/cloud/orchestration/events"
	"atonixcorp/cloud/orchestration/scheduler"
	"atonixcorp/cloud/orchestration/state"
)

func Create(id string, name string, plan string) error {
	sched := scheduler.NewScheduler()
	st := state.NewStateManager()
	ev := events.NewRouter()

	// Desired state
	st.SetDesired(id, "created")

	// Pick node for DB placement
	node, err := sched.PickNode(4, 8192)
	if err != nil {
		return err
	}

	// Backend call (placeholder)
	fmt.Println("Creating DB", name, "plan:", plan, "on node:", node)

	// Actual state
	st.SetActual(id, "created")

	// Emit event
	ev.Emit(events.Event{
		Type:      "db.create",
		Message:   fmt.Sprintf("DB %s created on %s", name, node),
		Timestamp: time.Now().Unix(),
	})

	return nil
}
