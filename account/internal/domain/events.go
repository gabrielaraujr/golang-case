package domain

import "github.com/google/uuid"

// Event types consumed by account service from risk-analysis.
// These constants define the contract between account and risk-analysis microservices.
//
// MESSAGE CONTRACT:
// - Flow: account (publishes ProposalCreated) → risk-analysis → account (consumes these events)
// - These event type values MUST match the constants defined in risk-analysis/internal/domain/events.go
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

// Event type published by account service to risk-analysis.
const (
	EventProposalCreated = "ProposalCreated"
)

// ProposalStatusChangedEvent represents an incoming event from risk-analysis service.
type ProposalStatusChangedEvent struct {
	EventType  string    `json:"event_type"`
	ProposalID uuid.UUID `json:"proposal_id"`
	Approved   bool      `json:"approved"`
}

type ProposalPayload struct {
	FullName string  `json:"full_name"`
	CPF      string  `json:"cpf"`
	Salary   float64 `json:"salary"`
}

// ProposalCreatedEvent represents an outgoing event to risk-analysis service.
// This must match the expected format in risk-analysis/internal/domain/events.go
type ProposalCreatedEvent struct {
	EventType  string           `json:"event_type"`
	ProposalID uuid.UUID        `json:"proposal_id"`
	Payload    *ProposalPayload `json:"payload"`
}
