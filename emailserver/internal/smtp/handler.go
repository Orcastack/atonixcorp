package smtp

// DefaultHandler is a no-op handler you can use if you don't need hooks yet.
type DefaultHandler struct{}

func (h *DefaultHandler) BeforeAccept(msg *Message) error {
	// Example:
	// - enforce max recipients
	// - block certain senders
	// - run extra validation
	return nil
}

func (h *DefaultHandler) AfterAccept(msg *Message) {
	// Example:
	// - send analytics event
	// - log delivery
	// - trigger webhooks
}
