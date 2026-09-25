package metrics

import "sync"

type Counter struct {
    mu    sync.Mutex
    value int64
}

func (c *Counter) Inc() {
    c.mu.Lock()
    c.value++
    c.mu.Unlock()
}

func (c *Counter) Value() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

var Requests = &Counter{}
var EventsIngested = &Counter{}
