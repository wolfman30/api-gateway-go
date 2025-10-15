package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/wolfman30/api-gateway-go/internal/models"
)

// CreateReel handles POST /reels
func CreateReel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateReelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate a unique run ID
	runID := uuid.New().String()

	// TODO: Publish command to SQS for orchestrator pickup
	log.Printf("Accepted reel request for project %s, runID=%s", req.ProjectID, runID)

	// Return 202 Accepted with runID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(models.CreateReelResponse{RunID: runID})
}

// GetRunStatus handles GET /runs/{runId}
func GetRunStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract runId from path
	runID := r.URL.Path[len("/runs/"):]
	if runID == "" {
		http.Error(w, "Missing runId", http.StatusBadRequest)
		return
	}

	// TODO: Query DynamoDB for run state
	log.Printf("Fetching status for runID=%s", runID)

	// Stub response
	resp := models.RunStatusResponse{
		RunID:  runID,
		Status: "PENDING",
		Steps:  []models.RunStep{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
