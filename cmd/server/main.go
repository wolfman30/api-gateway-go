package main

import (
	"log"
	"net/http"

	"github.com/wolfman30/api-gateway-go/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/reels", handlers.CreateReel)
	mux.HandleFunc("/runs/", handlers.GetRunStatus)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := ":8081"
	log.Printf("Starting API gateway on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
