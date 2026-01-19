package ports

import "context"

type ProposalEvent struct {
	EventType  string      `json:"event_type"`
	ProposalID string      `json:"proposal_id"`
	Payload    interface{} `json:"payload"`
}

type QueueProducer interface {
	Publish(ctx context.Context, event *ProposalEvent) error
}
