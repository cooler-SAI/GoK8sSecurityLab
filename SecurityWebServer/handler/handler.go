package handler

import (
	"fmt"
	"log"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Add security header: prevents clickjacking
	w.Header().Set("X-Frame-Options", "DENY")

	// Add Content-Type header for better security
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Output message that server is running
	_, err := fmt.Fprintf(w, "Status: OK. Running securely in a container.\n")
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
	log.Printf("Health check request received from %s", r.RemoteAddr)
}

// SecureHeaders Additional security headers for enhanced protection
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers to prevent various attacks
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}
