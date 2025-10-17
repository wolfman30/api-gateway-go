package main

import (
	"context"
	"log"
	"net/http"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/wolfman30/api-gateway-go/internal/bus"
	"github.com/wolfman30/api-gateway-go/internal/config"
	"github.com/wolfman30/api-gateway-go/internal/handlers"
)

func main() {
	ctx := context.Background()

	// Load secrets from AWS Secrets Manager
	secrets, err := config.LoadFromSecretsManager(ctx)
	if err != nil {
		log.Fatalf("Failed to load secrets: %v", err)
	}
	_ = secrets // Will be used for authentication/database connection

	// Load environment configuration
	envConfig := config.LoadEnvironmentConfig()
	log.Printf("Running in environment: %s", envConfig.Environment)

	// Load AWS configuration
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Create SQS client
	sqsClient := sqs.NewFromConfig(awsCfg)

	// Initialize SQS publisher with configured queue URL
	publisher := bus.NewPublisher(envConfig.SqsQueueURL, sqsClient)
	handlers.SetPublisher(publisher)

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/reels", handlers.CreateReel)
	mux.HandleFunc("/runs/", handlers.GetRunStatus)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := ":" + envConfig.ApiPort
	log.Printf("Starting API gateway on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
