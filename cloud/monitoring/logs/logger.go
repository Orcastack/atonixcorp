package logs

import "fmt"

type Logger struct {
	sinks []Sink
}

type Sink interface {
	Write(LogEntry) error
}

type LogEntry struct {
	Level   string
	Message string
	Fields  map[string]any
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) RegisterSink(s Sink) {
	l.sinks = append(l.sinks, s)
}

func (l *Logger) Log(level, msg string, fields map[string]any) {
	entry := LogEntry{Level: level, Message: msg, Fields: fields}
	for _, s := range l.sinks {
		_ = s.Write(entry)
	}
	fmt.Println("Monitoring: log:", level, msg)
}
