package domain

import (
	"errors"

	"github.com/google/uuid"
)

// Event types published by risk-analysis service after analyzing proposals.
// These constants define the contract between risk-analysis and account microservices.
//
// MESSAGE CONTRACT:
// - Flow: account → risk-analysis (consumes ProposalCreated) → account (publishes these events)
// - These event type values MUST match the constants defined in account/internal/domain/events.go
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

// ProposalStatusChangedEvent represents an outgoing event to account service.
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

// ProposalCreatedEvent represents an incoming event from account service.
// This must match the expected format in account/internal/domain/events.go
type ProposalCreatedEvent struct {
	EventType  string           `json:"event_type"`
	ProposalID uuid.UUID        `json:"proposal_id"`
	Payload    *ProposalPayload `json:"payload"`
}

var (
	ErrEmptyCPF       = errors.New("cpf is required")
	ErrNilPayload     = errors.New("payload cannot be nil")
	ErrEmptyFullName  = errors.New("full_name is required")
	ErrNilProposalID  = errors.New("proposal_id cannot be nil")
	ErrEmptyEventType = errors.New("event_type is required")
	ErrNegativeSalary = errors.New("salary cannot be negative")
)

func (e *ProposalCreatedEvent) Validate() error {
	if e.EventType == "" {
		return ErrEmptyEventType
	}
	if e.ProposalID == uuid.Nil {
		return ErrNilProposalID
	}
	if e.Payload == nil {
		return ErrNilPayload
	}
	if e.Payload.FullName == "" {
		return ErrEmptyFullName
	}
	if e.Payload.CPF == "" {
		return ErrEmptyCPF
	}
	if e.Payload.Salary < 0 {
		return ErrNegativeSalary
	}
	return nil
}
