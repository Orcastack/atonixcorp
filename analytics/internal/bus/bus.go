package bus

import (
    "sync"

    "atonixcorp/analytics/internal/events"
)

type Bus struct {
    mu       sync.RWMutex
    handlers []func(events.Event)
}

func New() *Bus {
    return &Bus{handlers: []func(events.Event){}}
}

func (b *Bus) Publish(evt events.Event) {
    b.mu.RLock()
    defer b.mu.RUnlock()
    for _, h := range b.handlers {
        go h(evt)
    }
}

func (b *Bus) Subscribe(handler func(events.Event)) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.handlers = append(b.handlers, handler)
}
