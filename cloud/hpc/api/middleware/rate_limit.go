package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]int
	reset    time.Time
	limit    int
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]int),
		reset:    time.Now().Add(1 * time.Minute),
		limit:    limit,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Context().Value("auth")
		if auth == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tenant := auth.(*AuthContext).TenantID

		rl.mu.Lock()

		if time.Now().After(rl.reset) {
			rl.requests = make(map[string]int)
			rl.reset = time.Now().Add(1 * time.Minute)
		}

		rl.requests[tenant]++

		if rl.requests[tenant] > rl.limit {
			rl.mu.Unlock()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
