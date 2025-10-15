package bus

import (
	"encoding/json"
	"log"
)

// Publisher handles publishing commands to SQS.
type Publisher struct {
	queueURL string
}

// NewPublisher creates a new SQS command publisher.
func NewPublisher(queueURL string) *Publisher {
	return &Publisher{queueURL: queueURL}
}

// PublishReelCommand sends a reel command to SQS for orchestrator pickup.
func (p *Publisher) PublishReelCommand(runID string, payload interface{}) error {
	// TODO: Use AWS SDK to send message to SQS
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Printf("Publishing reel command for runID=%s to queue=%s (stub): %s", runID, p.queueURL, string(body))
	return nil
}
