package main

import (
	"SecurityWebServer/handler"
	"SecurityWebServer/middleware"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Middleware chains for different endpoints
	healthChain := handler.RateLimit(
		handler.SecureHeaders(
			http.HandlerFunc(handler.HealthHandler),
		),
	)

	halloweenChain := handler.RateLimit(
		handler.SecureHeaders(
			http.HandlerFunc(handler.HalloweenHandler),
		),
	)

	apiChain := handler.RateLimit(
		handler.SecureHeaders(
			http.HandlerFunc(handler.HalloweenAPIHandler),
		),
	)

	infoChain := handler.RateLimit(
		handler.SecureHeaders(
			http.HandlerFunc(handler.InfoHandler),
		),
	)

	// Secure greet endpoint
	secureGreetChain := handler.RateLimit(
		handler.SecureHeaders(
			http.HandlerFunc(handler.SecureGreetHandler),
		),
	)

	// Vulnerable greet endpoint (for demonstration)
	vulnerableGreetChain := handler.RateLimit(
		handler.SecureHeaders(
			http.HandlerFunc(handler.VulnerableGreetHandler),
		),
	)

	// Register routes
	http.Handle("/", healthChain)
	http.Handle("/halloween", halloweenChain)
	http.Handle("/api/halloween", apiChain)
	http.Handle("/info", infoChain)
	// Secure greet endpoint
	http.Handle("/greet", secureGreetChain)
	// Vulnerable greet endpoint (for educational purposes)
	http.Handle("/vulnerable-greet", vulnerableGreetChain)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Printf("🎃 Starting Secure Halloween Server on port %s", port)
	log.Printf("📋 Available endpoints:")
	log.Printf("   🌐 Health Check:        http://localhost:%s/", port)
	log.Printf("   🎃 Halloween Page:      http://localhost:%s/halloween", port)
	log.Printf("   🔗 Halloween API:       http://localhost:%s/api/halloween", port)
	log.Printf("   ℹ️  Server Info:         http://localhost:%s/info", port)
	log.Printf("   ✅ SECURE GREET:         http://localhost:%s/greet?name=YourName", port)
	log.Printf("   ⚠️  VULNERABLE GREET:    http://localhost:%s/vulnerable-greet?name=YourName", port)
	log.Printf("⚡ Rate Limiting: 10 requests/second")
	log.Printf("🔒 Security headers enabled")
	log.Printf("🛡️  XSS Protection: Secure handlers use html/template")

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error starting server: %v", err)
	}
}
