package main

import (
	"SecurityWebServer/handler"
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

	// Цепочки middleware для разных endpoints
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

	// Регистрация routes
	http.Handle("/", healthChain)
	http.Handle("/halloween", halloweenChain)
	http.Handle("/api/halloween", apiChain)
	http.Handle("/info", infoChain)

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
	log.Printf("   🌐 Health Check:    http://localhost:%s/", port)
	log.Printf("   🎃 Halloween Page:  http://localhost:%s/halloween", port)
	log.Printf("   🔗 Halloween API:   http://localhost:%s/api/halloween", port)
	log.Printf("   ℹ️  Server Info:     http://localhost:%s/info", port)
	log.Printf("⚡ Rate Limiting: 10 requests/second")
	log.Printf("🔒 Security headers enabled")

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error starting server: %v", err)
	}
}
