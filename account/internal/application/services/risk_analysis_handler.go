package services

import (
	"context"

	events "github.com/gabrielaraujr/golang-case/account/internal/domain"
	"github.com/gabrielaraujr/golang-case/account/internal/domain/entities"
	"github.com/gabrielaraujr/golang-case/account/internal/ports"
)

type ProposalStatusChangedEventHandler struct {
	repository ports.ProposalRepository
	logger     ports.Logger
}

func NewProposalStatusChangedEventHandler(
	repo ports.ProposalRepository,
	logger ports.Logger,
) *ProposalStatusChangedEventHandler {
	return &ProposalStatusChangedEventHandler{
		repository: repo,
		logger:     logger,
	}
}

func (h *ProposalStatusChangedEventHandler) Handle(
	ctx context.Context,
	event *events.ProposalStatusChangedEvent,
) error {
	h.logger.Info(ctx, "processing risk analysis event", "event_type", event.EventType, "proposal_id", event.ProposalID)

	proposal, err := h.repository.FindByID(ctx, event.ProposalID)
	if err != nil {
		h.logger.Error(ctx, "proposal not found", "proposal_id", event.ProposalID, "error", err)
		return err
	}

	switch event.EventType {
	case events.EventDocumentsApproved:
		return h.handleAnalyzing(ctx, proposal)
	case events.EventDocumentsRejected, events.EventCreditRejected, events.EventFraudRejected:
		return h.handleRejection(ctx, proposal)
	case events.EventRiskAnalysisCompleted:
		return h.handleCompletion(ctx, proposal, event)
	default:
		h.logger.Info(ctx, "intermediate event received", "event_type", event.EventType)
		return nil
	}
}

func (h *ProposalStatusChangedEventHandler) handleAnalyzing(ctx context.Context, proposal *entities.Proposal) error {
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

func (h *ProposalStatusChangedEventHandler) handleRejection(ctx context.Context, proposal *entities.Proposal) error {
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

func (h *ProposalStatusChangedEventHandler) handleCompletion(
	ctx context.Context,
	proposal *entities.Proposal,
	event *events.ProposalStatusChangedEvent,
) error {
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
