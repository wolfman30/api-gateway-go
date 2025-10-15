package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/wolfman30/api-gateway-go/internal/bus"
	"github.com/wolfman30/api-gateway-go/internal/handlers"
)

func main() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Create SQS client
	sqsClient := sqs.NewFromConfig(cfg)

	// Initialize SQS publisher
	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		queueURL = "https://sqs.us-east-1.amazonaws.com/123456789012/reel-commands" // stub
	}
	publisher := bus.NewPublisher(queueURL, sqsClient)
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
