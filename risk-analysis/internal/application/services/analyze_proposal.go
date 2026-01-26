package services

import (
	"context"

	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain"
	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/ports"
	"github.com/google/uuid"
)

type AnalyzeProposalService struct {
	producer ports.QueueProducer
	logger   ports.Logger
}

func NewAnalyzeProposalService(producer ports.QueueProducer, logger ports.Logger) *AnalyzeProposalService {
	return &AnalyzeProposalService{
		producer: producer,
		logger:   logger,
	}
}

func (s *AnalyzeProposalService) Handle(ctx context.Context, event *domain.ProposalCreatedEvent) error {
	payload := event.Payload
	proposalID := event.ProposalID

	s.logger.Info(ctx, "[RiskAnalysis] Processing proposal", "proposal_id", proposalID)

	if err := event.Validate(); err != nil {
		s.logger.Error(ctx, "[RiskAnalysis] Invalid payload", "proposal_id", proposalID, "error", err)
		return err
	}

	// Document analysis
	documentResult := domain.AnalyzeDocuments(payload)
	documentApproved := documentResult.Approved
	if !documentApproved {
		s.logger.Warn(ctx, "[RiskAnalysis] Documents rejected", "proposal_id", proposalID, "reason", documentResult.Reason)
		return s.publish(ctx, domain.EventDocumentsRejected, proposalID, documentApproved)
	}

	s.logger.Info(ctx, "[RiskAnalysis] Documents approved", "proposal_id", proposalID)
	if err := s.publish(ctx, domain.EventDocumentsApproved, proposalID, documentApproved); err != nil {
		return err
	}

	// Credit analysis
	creditResult := domain.AnalyzeCredit(payload)
	if !creditResult.Approved {
		s.logger.Warn(ctx, "[RiskAnalysis] Credit rejected", "proposal_id", proposalID, "reason", creditResult.Reason)
		return s.publish(ctx, domain.EventCreditRejected, proposalID, creditResult.Approved)
	}

	// Fraud analysis
	fraudResult := domain.AnalyzeFraud(payload)
	if !fraudResult.Approved {
		s.logger.Warn(ctx, "[RiskAnalysis] Fraud rejected", "proposal_id", proposalID, "reason", fraudResult.Reason)
		return s.publish(ctx, domain.EventFraudRejected, proposalID, fraudResult.Approved)
	}

	// All analyses passed
	s.logger.Info(ctx, "[RiskAnalysis] Proposal fully approved", "proposal_id", proposalID)
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
