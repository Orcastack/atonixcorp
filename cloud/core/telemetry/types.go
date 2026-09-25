package telemetry

type Metric struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type Event struct {
	EventType string `json:"event_type"`
	Timestamp string `json:"timestamp"`
}
