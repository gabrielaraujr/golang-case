package services

import (
	"context"
	"log"

	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain/entities"
	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain/events"
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

func (s *AnalyzeProposalService) Handle(ctx context.Context, event *entities.IncomingEvent) error {
	log.Printf("[RiskAnalysis] Processing proposal %s", event.ProposalID)

	payload, err := event.ParsePayload()
	if err != nil {
		log.Printf("[RiskAnalysis] Failed to parse payload: %v", err)
		return err
	}

	// Document Analyze
	if !s.analyzeDocuments(payload) {
		log.Printf("[RiskAnalysis] Documents rejected for proposal %s", event.ProposalID)
		return s.publishResult(ctx, event.ProposalID, events.EventDocumentsApproved, events.EventDocumentsRejected, false)
	}

	if err := s.publishResult(
		ctx,
		event.ProposalID,
		events.EventDocumentsApproved,
		events.EventDocumentsRejected,
		true,
	); err != nil {
		return err
	}

	// Credit Analyze
	if !s.analyzeCredit(payload) {
		log.Printf("[RiskAnalysis] Credit rejected for proposal %s", event.ProposalID)
		return s.publishResult(ctx, event.ProposalID, events.EventCreditApproved, events.EventCreditRejected, false)
	}

	// Fraud Analyze
	if !s.analyzeFraud(payload) {
		log.Printf("[RiskAnalysis] Fraud rejected for proposal %s", event.ProposalID)
		return s.publishResult(ctx, event.ProposalID, events.EventFraudApproved, events.EventFraudRejected, false)
	}

	// All approved
	log.Printf("[RiskAnalysis] Proposal %s fully approved", event.ProposalID)
	return s.publishFinalResult(ctx, event.ProposalID, true)
}

func (s *AnalyzeProposalService) analyzeDocuments(payload *entities.ProposalPayload) bool {
	if len(payload.CPF) != 11 {
		log.Printf("[RiskAnalysis] Documents rejected: invalid CPF length %d", len(payload.CPF))
		return false
	}
	if len(payload.FullName) < 3 {
		log.Printf("[RiskAnalysis] Documents rejected: full name too short %s", payload.FullName)
		return false
	}

	log.Printf("[RiskAnalysis] Documents approved for CPF %s", payload.CPF)
	return true
}

func (s *AnalyzeProposalService) analyzeCredit(payload *entities.ProposalPayload) bool {
	if payload.Salary <= 3000.0 {
		log.Printf("[RiskAnalysis] Credit analysis rejected: insufficient salary %.2f", payload.Salary)
		return false
	}

	log.Printf("[RiskAnalysis] Credit analysis approved for CPF %s", payload.CPF)
	return true
}

func (s *AnalyzeProposalService) analyzeFraud(payload *entities.ProposalPayload) bool {
	lastDigit := payload.CPF[len(payload.CPF)-1] - '0'
	if lastDigit%2 != 0 {
		log.Printf("[RiskAnalysis] Fraud analysis rejected for CPF %s", payload.CPF)
		return false
	}

	log.Printf("[RiskAnalysis] Fraud analysis approved for CPF %s", payload.CPF)
	return true
}

func (s *AnalyzeProposalService) publishResult(
	ctx context.Context,
	proposalID uuid.UUID,
	approvedType,
	rejectedType string,
	approved bool,
) error {
	eventType := approvedType
	if !approved {
		eventType = rejectedType
	}
	return s.producer.Publish(ctx, &events.RiskEvent{
		EventType:  eventType,
		ProposalID: proposalID,
		Approved:   approved,
	})
}

func (s *AnalyzeProposalService) publishFinalResult(ctx context.Context, proposalID uuid.UUID, approved bool) error {
	return s.producer.Publish(ctx, &events.RiskEvent{
		EventType:  events.EventRiskAnalysisCompleted,
		ProposalID: proposalID,
		Approved:   approved,
	})
}
