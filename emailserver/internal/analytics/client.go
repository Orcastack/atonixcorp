package analytics

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Logger interface {
	Info(...any)
	Error(...any)
}

// ─────────────────────────────────────────────
// Event model
// ─────────────────────────────────────────────

type Event struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	UserID    string                 `json:"user_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ─────────────────────────────────────────────
// Transport interface (HTTP, Kafka, NATS, etc.)
// ─────────────────────────────────────────────

type Transport interface {
	Send(events []Event) error
}

// ─────────────────────────────────────────────
// HTTP transport (default)
// ─────────────────────────────────────────────

type HTTPTransport struct {
	Endpoint string
	Client   *http.Client
}

func NewHTTPTransport(endpoint string) *HTTPTransport {
	return &HTTPTransport{
		Endpoint: endpoint,
		Client:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (t *HTTPTransport) Send(events []Event) error {
	data, _ := json.Marshal(events)
	req, err := http.NewRequest("POST", t.Endpoint, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = t.Client.Do(req)
	return err
}

// ─────────────────────────────────────────────
// Analytics client
// ─────────────────────────────────────────────

type Client struct {
	transport Transport
	logger    Logger

	mu     sync.Mutex
	buffer []Event

	maxBatch int
	interval time.Duration

	stop chan struct{}
}

func NewClient(transport Transport, logger Logger) *Client {
	c := &Client{
		transport: transport,
		logger:    logger,
		maxBatch:  50,
		interval:  2 * time.Second,
		stop:      make(chan struct{}),
	}

	go c.loop()
	return c
}

// Fire-and-forget event
func (c *Client) Track(eventType string, userID string, meta map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buffer = append(c.buffer, Event{
		Type:      eventType,
		Timestamp: time.Now(),
		UserID:    userID,
		Metadata:  meta,
	})
}

// Background dispatcher
func (c *Client) loop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.flush()
		case <-c.stop:
			c.flush()
			return
		}
	}
}

func (c *Client) flush() {
	c.mu.Lock()
	if len(c.buffer) == 0 {
		c.mu.Unlock()
		return
	}

	// batch
	var batch []Event
	if len(c.buffer) > c.maxBatch {
		batch = c.buffer[:c.maxBatch]
		c.buffer = c.buffer[c.maxBatch:]
	} else {
		batch = c.buffer
		c.buffer = nil
	}
	c.mu.Unlock()

	// send
	if err := c.transport.Send(batch); err != nil {
		c.logger.Error("Analytics send failed:", err)
		// retry by putting back into buffer
		c.mu.Lock()
		c.buffer = append(batch, c.buffer...)
		c.mu.Unlock()
		return
	}

	c.logger.Info("Analytics sent batch:", len(batch))
}

func (c *Client) Close() {
	close(c.stop)
}
