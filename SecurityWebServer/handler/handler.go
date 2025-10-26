package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

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

	now := time.Now()
	if last, exists := rl.ips[ip]; exists {
		if now.Sub(last) < rl.window {
			return false
		}
	}

	rl.ips[ip] = now
	return true
}

// RateLimit middleware
func RateLimit(next http.Handler) http.Handler {
	limiter := NewRateLimiter(100 * time.Millisecond)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if !limiter.Allow(ip) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecureHeaders middleware
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

// HealthHandler - базовый health check
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := fmt.Fprintf(w, "Status: OK. Running securely in a container.\n")
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
	log.Printf("Health check request received from %s", r.RemoteAddr)
}
