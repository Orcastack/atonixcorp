package audit

// Event represents a structured audit event.
type Event struct {
	Name   string
	Fields map[string]any
}

// Common platform audit events
var (
	EventLogin            = Event{Name: "login"}
	EventLogout           = Event{Name: "logout"}
	EventDataAccess       = Event{Name: "data_access"}
	EventDataWrite        = Event{Name: "data_write"}
	EventProcessStart     = Event{Name: "process_start"}
	EventProcessComplete  = Event{Name: "process_complete"}
	EventSecurityAlert    = Event{Name: "security_alert"}
	EventAnomalyDetected  = Event{Name: "anomaly_detected"}
	EventPermissionDenied = Event{Name: "permission_denied"}
)
