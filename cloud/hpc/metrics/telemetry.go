package metrics

import (
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryConfig struct {
	ServiceName string
}

type Telemetry struct {
	Tracer trace.Tracer
}

func NewTelemetry(cfg TelemetryConfig) *Telemetry {
	// In production, plug in OTLP exporter + resource attributes.
	// Here we just use the global tracer provider.
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer(cfg.ServiceName)

	log.Printf("[telemetry] initialized for service=%s at=%s", cfg.ServiceName, time.Now().UTC())

	return &Telemetry{
		Tracer: tracer,
	}
}
