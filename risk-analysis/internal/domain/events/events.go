package events

import "github.com/google/uuid"

// Event types published by risk-analysis service after analyzing proposals.
// These constants define the contract between risk-analysis and account microservices.
//
// MESSAGE CONTRACT:
// - Flow: account → risk-analysis (consumes ProposalCreated) → account (publishes these events)
// - These event type values MUST match the constants defined in account/internal/ports/queue_consumer.go
// - Any changes to these values require coordinated deployment of both services
//
// Validation Flow:
//  1. Documents: CPF length (11) and full name length (≥3)
//  2. Credit: Salary threshold (>3000)
//  3. Fraud: CPF last digit parity check (even = approved)
//  4. RiskAnalysisCompleted: Published when all validations pass
const (
	EventDocumentsApproved     = "DocumentsApproved"
	EventDocumentsRejected     = "DocumentsRejected"
	EventCreditApproved        = "CreditApproved"
	EventCreditRejected        = "CreditRejected"
	EventFraudApproved         = "FraudApproved"
	EventFraudRejected         = "FraudRejected"
	EventRiskAnalysisCompleted = "RiskAnalysisCompleted"
)

// Event type consumed by risk-analysis service from account.
const (
	EventProposalCreated = "ProposalCreated"
)

type RiskEvent struct {
	EventType  string    `json:"event_type"`
	ProposalID uuid.UUID `json:"proposal_id"`
	Approved   bool      `json:"approved"`
}
