package events

import "github.com/google/uuid"

const (
	EventDocumentsApproved     = "DocumentsApproved"
	EventDocumentsRejected     = "DocumentsRejected"
	EventCreditApproved        = "CreditApproved"
	EventCreditRejected        = "CreditRejected"
	EventFraudApproved         = "FraudApproved"
	EventFraudRejected         = "FraudRejected"
	EventRiskAnalysisCompleted = "RiskAnalysisCompleted"
)

const (
	EventProposalCreated = "ProposalCreated"
)

type RiskEvent struct {
	EventType  string    `json:"event_type"`
	ProposalID uuid.UUID `json:"proposal_id"`
	Approved   bool      `json:"approved"`
}
