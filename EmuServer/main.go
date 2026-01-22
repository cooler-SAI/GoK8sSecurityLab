package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Terminal colors
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
	ColorCyan   = "\033[36m"
)

type Config struct {
	Port      string
	MaxDelay  int
	ErrorRate int
}

func main() {
	cfg := Config{}
	flag.StringVar(&cfg.Port, "port", "8080", "Port to run the server on")
	flag.IntVar(&cfg.MaxDelay, "delay", 0, "Maximum response delay in ms")
	flag.IntVar(&cfg.ErrorRate, "error-rate", 0, "Error chance percentage (0-100)")
	flag.Parse()

	mux := http.NewServeMux()

	// 1. Handlers are now clean - only business logic
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/health", handleHealth)

	// 2. Wrap everything in Middleware (Chain of Responsibility)
	// First log, then apply chaos, then execute logic
	finalHandler := loggerMiddleware(chaosMiddleware(cfg, mux))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      finalHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 3. Graceful Shutdown (Critical for K8s)
	// Don't let the server crash instantly, allow requests to finish processing
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("%s🚀 EmuServer started on :%s%s\n", ColorGreen, cfg.Port, ColorReset)
		fmt.Printf("%s⚙️  Chaos: delay=%dms, errors=%d%%%s\n", ColorCyan, cfg.MaxDelay, cfg.ErrorRate, ColorReset)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	<-done // Wait for signal (Ctrl+C or kill from K8s)
	fmt.Printf("\n%s🛑 Shutting down gracefully...%s\n", ColorYellow, ColorReset)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	fmt.Println("👋 Server exited")
}

// --- Middleware ---

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("[%s] %s%s%s %s (Lat: %s)\n",
			time.Now().Format("15:04:05"),
			ColorYellow, r.Method, ColorReset,
			r.URL.Path,
			time.Since(start),
		)
	})
}

func chaosMiddleware(cfg Config, next http.Handler) http.Handler {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip /health, otherwise K8s will kill the container for "unhealthiness"
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		// Simulate delay
		if cfg.MaxDelay > 0 {
			delay := rng.Intn(cfg.MaxDelay)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		// Simulate error
		if cfg.ErrorRate > 0 && rng.Intn(100) < cfg.ErrorRate {
			fmt.Printf("%s[CHAOS] 500 Error for %s%s\n", ColorRed, r.URL.Path, ColorReset)
			http.Error(w, "Chaos Monkey Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// --- Handlers ---

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprintf(w, "EmuServer: %s\nTime: %s", r.URL.Path, time.Now().Format(time.RFC3339))
	if err != nil {
		return
	}
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
