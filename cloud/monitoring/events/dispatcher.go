package events

import "fmt"

type Dispatcher struct {
	sinks []Sink
}

type Sink interface {
	Send(Event) error
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) RegisterSink(s Sink) {
	d.sinks = append(d.sinks, s)
}

func (d *Dispatcher) Dispatch(e Event) {
	for _, s := range d.sinks {
		_ = s.Send(e)
	}
	fmt.Println("Monitoring: event dispatched:", e.Type)
}
