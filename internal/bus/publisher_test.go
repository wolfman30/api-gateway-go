package bus

import (
	"testing"
)

func TestNewPublisher(t *testing.T) {
	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
	pub := NewPublisher(queueURL)

	if pub == nil {
		t.Fatal("Expected non-nil publisher")
	}

	if pub.queueURL != queueURL {
		t.Errorf("Expected queueURL %s, got %s", queueURL, pub.queueURL)
	}
}

func TestPublishReelCommand(t *testing.T) {
	pub := NewPublisher("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue")

	payload := map[string]string{
		"projectId": "proj_123",
		"idea":      "Test reel idea",
	}

	err := pub.PublishReelCommand("run-456", payload)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test with invalid payload that can't be marshaled
	invalidPayload := make(chan int) // channels can't be marshaled to JSON
	err = pub.PublishReelCommand("run-789", invalidPayload)
	if err == nil {
		t.Error("Expected error when marshaling invalid payload")
	}
}
