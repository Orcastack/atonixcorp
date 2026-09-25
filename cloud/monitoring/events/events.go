package events

type Event struct {
	Type      string
	Timestamp string
	Payload   map[string]any
}
