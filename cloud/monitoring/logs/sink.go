package logs

import (
	"fmt"
	"os"
)

// Sink is the interface all log destinations must implement.
type Sink interface {
	Write(entry LogEntry) error
}

// ------------------------------------------------------------
// 1. Stdout Sink
// ------------------------------------------------------------
type StdoutSink struct{}

func NewStdoutSink() *StdoutSink {
	return &StdoutSink{}
}

func (s *StdoutSink) Write(entry LogEntry) error {
	fmt.Printf("[LOG][%s] %s %+v\n", entry.Level, entry.Message, entry.Fields)
	return nil
}

// ------------------------------------------------------------
// 2. File Sink
// ------------------------------------------------------------
type FileSink struct {
	Path string
}

func NewFileSink(path string) *FileSink {
	return &FileSink{Path: path}
}

func (f *FileSink) Write(entry LogEntry) error {
	line := fmt.Sprintf("[%s] %s %+v\n", entry.Level, entry.Message, entry.Fields)
	return os.WriteFile(f.Path, []byte(line), 0644)
}

// ------------------------------------------------------------
// 3. Placeholder for enterprise sinks
// ------------------------------------------------------------

// Swift / S3 / Object Storage
type ObjectStorageSink struct{}

func (o *ObjectStorageSink) Write(entry LogEntry) error {
	// TODO: integrate Swift or S3 client
	fmt.Println("ObjectStorageSink: write log (TODO)")
	return nil
}

// Elasticsearch / OpenSearch
type ElasticSink struct{}

func (e *ElasticSink) Write(entry LogEntry) error {
	// TODO: send JSON to Elasticsearch index
	fmt.Println("ElasticSink: write log (TODO)")
	return nil
}

// Loki
type LokiSink struct{}

func (l *LokiSink) Write(entry LogEntry) error {
	// TODO: push log entry to Loki
	fmt.Println("LokiSink: write log (TODO)")
	return nil
}
