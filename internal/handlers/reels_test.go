package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wolfman30/api-gateway-go/internal/models"
)

func TestCreateReel(t *testing.T) {
	// Sample payload matching the digital marketing ICP example
	payload := models.CreateReelRequest{
		ProjectID: "proj_789",
		ICP: models.IdealClientProfile{
			Industry:           "Digital marketing",
			AudiencePainPoints: []string{"Creating consistent content takes too much time"},
			DesiredOutcome:     "Effortlessly generate high-quality AI twin reels to scale content production",
		},
		Idea: "Show how AI twins let you create reels in minutes instead of hours",
		FluxModel: models.FluxModelConfig{
			LoraURL:  "https://v3.fal.media/files/elephant/T6tBgeMb8efOTD9xv2cif_pytorch_lora_weights.safetensors",
			CfgScale: 8,
			Steps:    30,
		},
		FluxPrompt: models.FluxPromptRequest{
			Prompt:         "Professional digital marketer in modern home office setup",
			NegativePrompt: "blurry, low-resolution",
			AspectRatio:    "9:16",
			BatchSize:      4,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/reels", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	CreateReel(rec, req)

	// Assert status code
	if rec.Code != http.StatusAccepted {
		t.Errorf("Expected status %d, got %d", http.StatusAccepted, rec.Code)
	}

	// Parse response
	var resp models.CreateReelResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify runID is non-empty
	if resp.RunID == "" {
		t.Error("Expected non-empty runID in response")
	}

	t.Logf("CreateReel returned runID: %s", resp.RunID)
}

func TestCreateReel_InvalidPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/reels", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	CreateReel(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateReel_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/reels", nil)
	rec := httptest.NewRecorder()

	CreateReel(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d for wrong method, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestGetRunStatus(t *testing.T) {
	runID := "test-run-123"
	req := httptest.NewRequest(http.MethodGet, "/runs/"+runID, nil)
	rec := httptest.NewRecorder()

	GetRunStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp models.RunStatusResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.RunID != runID {
		t.Errorf("Expected runID %s, got %s", runID, resp.RunID)
	}

	if resp.Status != "PENDING" {
		t.Errorf("Expected status PENDING, got %s", resp.Status)
	}
}

func TestGetRunStatus_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/runs/test-run", nil)
	rec := httptest.NewRecorder()

	GetRunStatus(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d for wrong method, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestGetRunStatus_MissingRunID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/runs/", nil)
	rec := httptest.NewRecorder()

	GetRunStatus(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for missing runID, got %d", http.StatusBadRequest, rec.Code)
	}
}
