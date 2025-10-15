package bus

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQSClient defines the interface for SQS operations (for testing).
type SQSClient interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

// Publisher handles publishing commands to SQS.
type Publisher struct {
	queueURL  string
	sqsClient SQSClient
}

// NewPublisher creates a new SQS command publisher.
func NewPublisher(queueURL string, sqsClient SQSClient) *Publisher {
	return &Publisher{
		queueURL:  queueURL,
		sqsClient: sqsClient,
	}
}

// PublishReelCommand sends a reel command to SQS for orchestrator pickup.
func (p *Publisher) PublishReelCommand(runID string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create the SQS message with the runID as a message attribute
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(body)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"runId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(runID),
			},
		},
	}

	_, err = p.sqsClient.SendMessage(context.Background(), input)
	if err != nil {
		log.Printf("Failed to send message to SQS for runID=%s: %v", runID, err)
		return err
	}

	log.Printf("Successfully published reel command for runID=%s to queue=%s", runID, p.queueURL)
	return nil
}
