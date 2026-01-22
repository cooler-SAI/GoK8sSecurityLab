package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Colors for beautiful terminal output
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
	ColorCyan   = "\033[36m"
)

func main() {
	// 1. Setup startup parameters (flags)
	port := flag.String("port", "8080", "Port to run the server on")
	latency := flag.Int("delay", 0, "Maximum response delay in milliseconds")
	errorRate := flag.Int("error-rate", 0, "Error chance percentage (0-100)")
	flag.Parse()

	// 2. Define routes (Handlers)
	mux := http.NewServeMux()

	// Main page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		applyChaos(w, *latency, *errorRate)
		_, err := fmt.Fprintf(w, "Server emulator running!\nPath: %s\nTime: %s", r.URL.Path, time.Now().Format(time.RFC1123))
		if err != nil {
			return
		}
	})

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	})

	// 3. Start server
	serverAddr := ":" + *port
	fmt.Printf("%s🚀 Emulator running on http://localhost%s%s\n", ColorGreen, serverAddr, ColorReset)
	fmt.Printf("%s⚙️  Configuration: delay up to %dms, error chance %d%%%s\n", ColorCyan, *latency, *errorRate, ColorReset)
	fmt.Println("--------------------------------------------------")

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Startup error: %v", err)
	}
}

// logRequest prints information about incoming request to terminal
func logRequest(r *http.Request) {
	fmt.Printf("[%s] %s %s %s %s\n",
		time.Now().Format("15:04:05"),
		ColorYellow, r.Method, ColorReset,
		r.URL.Path)
}

// applyChaos adds delays and random errors
func applyChaos(w http.ResponseWriter, maxDelay int, errorRate int) {
	// Create a new local random generator for thread safety
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Simulate delay
	if maxDelay > 0 {
		delay := rng.Intn(maxDelay)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	// Simulate error
	if errorRate > 0 && rng.Intn(100) < errorRate {
		fmt.Printf("%s[CHAOS] Generated random error 500%s\n", ColorRed, ColorReset)
		http.Error(w, "Internal Server Error (Chaos Monkey)", http.StatusInternalServerError)
		return
	}
}
