package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct {
	Collectors *HPCCollectors
	Registry   *prometheus.Registry
}

func NewExporter() *Exporter {
	reg := prometheus.NewRegistry()
	collectors := NewHPCCollectors()
	collectors.Register(reg)

	return &Exporter{
		Collectors: collectors,
		Registry:   reg,
	}
}

func (e *Exporter) Handler() http.Handler {
	return promhttp.HandlerFor(e.Registry, promhttp.HandlerOpts{})
}

func (e *Exporter) Serve(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", e.Handler())

	log.Printf("[metrics] exporter listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}
