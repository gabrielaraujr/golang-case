package services

import (
	"context"
	"log"

	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain"
	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/ports"
	"github.com/google/uuid"
)

type AnalyzeProposalService struct {
	producer ports.QueueProducer
}

func NewAnalyzeProposalService(producer ports.QueueProducer) *AnalyzeProposalService {
	return &AnalyzeProposalService{
		producer: producer,
	}
}

func (s *AnalyzeProposalService) Handle(ctx context.Context, event *domain.ProposalCreatedEvent) error {
	log.Printf("[RiskAnalysis] Processing proposal %s", event.ProposalID)

	if err := event.Validate(); err != nil {
		log.Printf("[RiskAnalysis] Invalid payload for proposal %s: %v", event.ProposalID, err)
		return err
	}

	payload := event.Payload
	proposalID := event.ProposalID

	// Document analysis
	documentResult := domain.AnalyzeDocuments(payload)
	documentApproved := documentResult.Approved
	if !documentApproved {
		log.Printf("[RiskAnalysis] Documents rejected for proposal %s: %s", proposalID, documentResult.Reason)
		return s.publish(ctx, domain.EventDocumentsRejected, proposalID, documentApproved)
	}

	log.Printf("[RiskAnalysis] Documents approved for proposal %s", proposalID)
	if err := s.publish(ctx, domain.EventDocumentsApproved, proposalID, documentApproved); err != nil {
		return err
	}

	// Credit analysis
	creditResult := domain.AnalyzeCredit(payload)
	if !creditResult.Approved {
		log.Printf("[RiskAnalysis] Credit rejected for proposal %s: %s", proposalID, creditResult.Reason)
		return s.publish(ctx, domain.EventCreditRejected, proposalID, creditResult.Approved)
	}

	// Fraud analysis
	fraudResult := domain.AnalyzeFraud(payload)
	if !fraudResult.Approved {
		log.Printf("[RiskAnalysis] Fraud rejected for proposal %s: %s", proposalID, fraudResult.Reason)
		return s.publish(ctx, domain.EventFraudRejected, proposalID, fraudResult.Approved)
	}

	// All analyses passed
	log.Printf("[RiskAnalysis] Proposal %s fully approved", proposalID)
	return s.publish(ctx, domain.EventRiskAnalysisCompleted, proposalID, true)
}

func (s *AnalyzeProposalService) publish(
	ctx context.Context,
	eventType string,
	proposalID uuid.UUID,
	approved bool,
) error {
	return s.producer.Publish(ctx, &domain.ProposalStatusChangedEvent{
		EventType:  eventType,
		ProposalID: proposalID,
		Approved:   approved,
	})
}
