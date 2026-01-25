package ports

import "context"

// Event types for risk analysis responses received from risk-analysis service.
// These constants define the contract between account and risk-analysis microservices.
//
// MESSAGE CONTRACT:
// - Flow: account (publishes ProposalCreated) → risk-analysis → account (consumes these events)
// - These event type values MUST match the constants defined in risk-analysis/internal/domain/events/events.go
// - Any changes to these values require coordinated deployment of both services
//
// Events:
//   - DocumentsApproved/Rejected: First validation step (CPF and name validation)
//   - CreditApproved/Rejected: Second validation step (salary threshold check)
//   - FraudApproved/Rejected: Third validation step (CPF last digit check)
//   - RiskAnalysisCompleted: Final result when all validations pass
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
