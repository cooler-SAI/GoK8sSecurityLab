package main

import (
	"SecurityWebServer/handler"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configure routing with security middleware
	securedHandler := handler.SecureHeaders(http.HandlerFunc(handler.HealthHandler))
	http.Handle("/", securedHandler)

	// Create custom server with timeouts for security
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		// Protection against slow DoS attacks
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	fmt.Printf("Starting server on port %s\n", port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
