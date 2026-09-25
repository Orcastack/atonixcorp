package detect

import (
	"time"
)

// Anomaly represents a detected suspicious event.
type Anomaly struct {
	Identity  string
	Event     string
	Details   map[string]any
	Timestamp time.Time
}

// Detector defines the interface for anomaly detection engines.
type Detector interface {
	Check(identity string, event string, fields map[string]any) (*Anomaly, error)
}

// SimpleDetector is a lightweight anomaly engine for the entire platform.
type SimpleDetector struct {
	rules []Rule
}

// NewSimpleDetector creates a detector with a set of rules.
func NewSimpleDetector(rules []Rule) *SimpleDetector {
	return &SimpleDetector{rules: rules}
}

// Check evaluates all rules against an event.
func (d *SimpleDetector) Check(identity string, event string, fields map[string]any) (*Anomaly, error) {
	for _, rule := range d.rules {
		if rule.Match(identity, event, fields) {
			return &Anomaly{
				Identity:  identity,
				Event:     event,
				Details:   fields,
				Timestamp: time.Now().UTC(),
			}, nil
		}
	}
	return nil, nil
}
