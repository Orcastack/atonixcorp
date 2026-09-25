package policy

import (
	"atonixcorp/cloud/controlplane/bus"
	"atonixcorp/cloud/controlplane/models"
)

type Engine struct {
	bus *bus.Bus
}

func NewEngine(b *bus.Bus) *Engine {
	e := &Engine{bus: b}
	b.Subscribe("workload.submitted", e.onWorkloadSubmitted)
	return e
}

func (e *Engine) onWorkloadSubmitted(ev bus.Event) {
	_, ok := ev.Data.(models.Workload)
	if !ok {
		return
	}
	// placeholder: always allow
	e.bus.Publish(bus.Event{
		Type: "policy.evaluated",
		Data: ev.Data,
	})
}
