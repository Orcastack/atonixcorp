package audit

import (
	"fmt"

	"atonixcorp/cloud/controlplane/bus"
)

type Service struct {
	bus *bus.Bus
}

func NewService(b *bus.Bus) *Service {
	s := &Service{bus: b}
	b.Subscribe("workload.submitted", s.onEvent)
	b.Subscribe("workload.scheduled", s.onEvent)
	b.Subscribe("policy.evaluated", s.onEvent)
	return s
}

func (s *Service) onEvent(ev bus.Event) {
	fmt.Println("Audit:", ev.Type)
}
