package handlers

import (
	"net/http"

	"atonixcorp/cloud/hpc/metrics"
)

type MetricsHandler struct {
	Exporter *metrics.Exporter
}

func NewMetricsHandler(e *metrics.Exporter) *MetricsHandler {
	return &MetricsHandler{Exporter: e}
}

func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Exporter.Handler().ServeHTTP(w, r)
}
