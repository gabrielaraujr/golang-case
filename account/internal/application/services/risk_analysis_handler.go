package services

import (
	"context"

	"github.com/gabrielaraujr/golang-case/account/internal/domain/entities"
	"github.com/gabrielaraujr/golang-case/account/internal/ports"
	"github.com/google/uuid"
)

type RiskAnalysisEventHandler struct {
	repository ports.ProposalRepository
	logger     ports.Logger
}

func NewRiskAnalysisEventHandler(repo ports.ProposalRepository, logger ports.Logger) *RiskAnalysisEventHandler {
	return &RiskAnalysisEventHandler{
		repository: repo,
		logger:     logger,
	}
}

func (h *RiskAnalysisEventHandler) Handle(ctx context.Context, event *ports.RiskAnalysisEvent) error {
	h.logger.Info(ctx, "processing risk analysis event", "event_type", event.EventType, "proposal_id", event.ProposalID)

	proposalID, err := uuid.Parse(event.ProposalID)
	if err != nil {
		h.logger.Error(ctx, "invalid proposal ID", "error", err)
		return err
	}

	proposal, err := h.repository.FindByID(ctx, proposalID)
	if err != nil {
		h.logger.Error(ctx, "proposal not found", "proposal_id", event.ProposalID, "error", err)
		return err
	}

	switch event.EventType {
	case ports.EventDocumentsApproved:
		return h.handleAnalyzing(ctx, proposal)
	case ports.EventDocumentsRejected, ports.EventCreditRejected, ports.EventFraudRejected:
		return h.handleRejection(ctx, proposal)
	case ports.EventRiskAnalysisCompleted:
		return h.handleCompletion(ctx, proposal, event)
	default:
		h.logger.Info(ctx, "intermediate event received", "event_type", event.EventType)
		return nil
	}
}

func (h *RiskAnalysisEventHandler) handleAnalyzing(ctx context.Context, proposal *entities.Proposal) error {
	if !proposal.IsPending() {
		return nil
	}

	if err := proposal.StartAnalysis(); err != nil {
		h.logger.Error(ctx, "failed to start analysis", "error", err)
		return err
	}

	if err := h.repository.Update(ctx, proposal); err != nil {
		h.logger.Error(ctx, "failed to update proposal", "error", err)
		return err
	}

	h.logger.Info(ctx, "proposal moved to analyzing", "proposal_id", proposal.ID.String())
	return nil
}

func (h *RiskAnalysisEventHandler) handleRejection(ctx context.Context, proposal *entities.Proposal) error {
	if err := proposal.Reject(); err != nil {
		h.logger.Error(ctx, "failed to reject proposal", "error", err)
		return err
	}

	if err := h.repository.Update(ctx, proposal); err != nil {
		h.logger.Error(ctx, "failed to update proposal", "error", err)
		return err
	}

	h.logger.Info(ctx, "proposal rejected", "proposal_id", proposal.ID.String())
	return nil
}

func (h *RiskAnalysisEventHandler) handleCompletion(ctx context.Context, proposal *entities.Proposal, event *ports.RiskAnalysisEvent) error {
	if !event.Approved {
		return h.handleRejection(ctx, proposal)
	}

	if err := proposal.Approve(); err != nil {
		h.logger.Error(ctx, "failed to approve proposal", "error", err)
		return err
	}

	if err := h.repository.Update(ctx, proposal); err != nil {
		h.logger.Error(ctx, "failed to update proposal", "error", err)
		return err
	}

	h.logger.Info(ctx, "proposal approved", "proposal_id", proposal.ID.String())
	return nil
}
