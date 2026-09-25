package detect

// Rule defines a security condition that triggers an anomaly.
type Rule interface {
	Match(identity string, event string, fields map[string]any) bool
}

// SuspiciousEventRule triggers on specific event names.
type SuspiciousEventRule struct {
	EventName string
}

func (r SuspiciousEventRule) Match(identity string, event string, fields map[string]any) bool {
	return event == r.EventName
}

// ThresholdRule triggers when a numeric field exceeds a threshold.
type ThresholdRule struct {
	Field     string
	Threshold float64
}

func (r ThresholdRule) Match(identity string, event string, fields map[string]any) bool {
	val, ok := fields[r.Field]
	if !ok {
		return false
	}

	num, ok := val.(float64)
	if !ok {
		return false
	}

	return num > r.Threshold
}

// IdentityRule triggers when a specific identity performs an event.
type IdentityRule struct {
	Identity string
	Event    string
}

func (r IdentityRule) Match(identity string, event string, fields map[string]any) bool {
	return identity == r.Identity && event == r.Event
}
