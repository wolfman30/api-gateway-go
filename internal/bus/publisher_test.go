package bus

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// MockSQSClient is a mock implementation of SQSClient for testing.
type MockSQSClient struct {
	SendMessageFunc func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

func (m *MockSQSClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(ctx, params, optFns...)
	}
	return &sqs.SendMessageOutput{}, nil
}

func TestNewPublisher(t *testing.T) {
	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
	mockClient := &MockSQSClient{}
	pub := NewPublisher(queueURL, mockClient)

	if pub == nil {
		t.Fatal("Expected non-nil publisher")
	}

	if pub.queueURL != queueURL {
		t.Errorf("Expected queueURL %s, got %s", queueURL, pub.queueURL)
	}

	if pub.sqsClient == nil {
		t.Error("Expected non-nil sqsClient")
	}
}

func TestPublishReelCommand(t *testing.T) {
	mockClient := &MockSQSClient{}
	pub := NewPublisher("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue", mockClient)

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
