package ports

import "context"

const (
	EventDocumentsApproved     = "DocumentsApproved"
	EventDocumentsRejected     = "DocumentsRejected"
	EventCreditApproved        = "CreditApproved"
	EventCreditRejected        = "CreditRejected"
	EventFraudApproved         = "FraudApproved"
	EventFraudRejected         = "FraudRejected"
	EventRiskAnalysisCompleted = "RiskAnalysisCompleted"
)

type RiskAnalysisEvent struct {
	EventType  string `json:"event_type"`
	ProposalID string `json:"proposal_id"`
	Approved   bool   `json:"approved"`
}

type EventHandler interface {
	Handle(ctx context.Context, event *RiskAnalysisEvent) error
}

type QueueMessage struct {
	ID         string
	Body       string
	Attributes map[string]string
}

type QueueConsumer interface {
	Start(ctx context.Context, handler func(context.Context, *QueueMessage) error) error
	Close() error
}
