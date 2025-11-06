package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// --- RATE LIMITER IMPLEMENTATION ---

// RateLimiter Simple structure for rate limiting
type RateLimiter struct {
	ips    map[string]time.Time
	mux    sync.Mutex
	window time.Duration
}

func NewRateLimiter(window time.Duration) *RateLimiter {
	return &RateLimiter{
		ips:    make(map[string]time.Time),
		window: window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mux.Lock()
	defer rl.mux.Unlock()

	// Window of 100ms means maximum 10 requests per second per IP
	if last, exists := rl.ips[ip]; exists {
		if time.Since(last) < rl.window {
			return false
		}
	}

	rl.ips[ip] = time.Now()
	return true
}

// --- MIDDLEWARES ---

// RateLimit middleware (EXPORTED - starts with capital letter)
func RateLimit(next http.Handler) http.Handler {
	// 100ms per request (10 requests per second)
	limiter := NewRateLimiter(100 * time.Millisecond)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if !limiter.Allow(ip) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests) // HTTP 429
			_, err := fmt.Fprintln(w, "429 Too Many Requests: Rate limit exceeded.")
			if err != nil {
				log.Printf("Error writing rate limit response: %v", err)
				return
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecureHeaders middleware (EXPORTED - starts with capital letter)
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

// HealthHandler - basic health check (EXPORTED)
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Status: OK. Running securely in a container.\n")
	if err != nil {
		log.Printf("Error writing health response: %v", err)
	}
}
