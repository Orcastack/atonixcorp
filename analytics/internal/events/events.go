package events

import "time"

type EventKind string

const (
    AuthEvent   EventKind = "auth"
    FileEvent   EventKind = "file"
    SystemEvent EventKind = "system"
    NodeEvent   EventKind = "node"
    DeviceEvent EventKind = "device"
    JobEvent    EventKind = "job"
)

type Event struct {
    ID        string                 `db:"id"`
    Kind      EventKind              `db:"kind"`
    Payload   map[string]any         `db:"payload"`
    Timestamp time.Time              `db:"ts"`
    Source    string                 `db:"source"`
}
