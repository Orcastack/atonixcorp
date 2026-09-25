package metrics

import "fmt"

type Exporter struct {
	sinks []Sink
}

type Sink interface {
	Push([]Metric) error
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) RegisterSink(s Sink) {
	e.sinks = append(e.sinks, s)
}

func (e *Exporter) Export(metrics []Metric) {
	for _, s := range e.sinks {
		_ = s.Push(metrics)
	}
	fmt.Println("Monitoring: metrics exported")
}
