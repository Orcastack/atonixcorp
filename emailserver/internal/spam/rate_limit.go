package spam

import (
	"sync"
	"time"
)

type RateLimiter interface {
	Allow(key string) bool
}

type MemoryRateLimiter struct {
	mu      sync.Mutex
	entries map[string]*entry
	limit   int           // max requests per window
	window  time.Duration // window size
}

type entry struct {
	count     int
	expiresAt time.Time
}

func NewMemoryRateLimiter(limit int, window time.Duration) *MemoryRateLimiter {
	return &MemoryRateLimiter{
		entries: make(map[string]*entry),
		limit:   limit,
		window:  window,
	}
}

func (r *MemoryRateLimiter) Allow(key string) bool {
	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	e, ok := r.entries[key]
	if !ok || now.After(e.expiresAt) {
		r.entries[key] = &entry{
			count:     1,
			expiresAt: now.Add(r.window),
		}
		return true
	}

	if e.count >= r.limit {
		return false
	}

	e.count++
	return true
}
