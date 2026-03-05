package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytesWritten, err := fmt.Fprintf(w, "hello webhook")
		if err != nil {
			log.Println("Error writing response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("Request received: %s, bytes written: %d", r.URL.Path, bytesWritten)
	})

	port := "8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
