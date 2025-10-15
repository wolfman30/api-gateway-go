package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wolfman30/api-gateway-go/internal/bus"
	"github.com/wolfman30/api-gateway-go/internal/handlers"
)

func main() {
	// Initialize SQS publisher
	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		queueURL = "https://sqs.us-east-1.amazonaws.com/123456789012/reel-commands" // stub
	}
	publisher := bus.NewPublisher(queueURL)
	handlers.SetPublisher(publisher)

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
