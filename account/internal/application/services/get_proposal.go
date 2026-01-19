package services

import (
	"context"
	"errors"

	appErrors "github.com/gabrielaraujr/golang-case/account/internal/application"
	"github.com/gabrielaraujr/golang-case/account/internal/application/dto"
	domainErrors "github.com/gabrielaraujr/golang-case/account/internal/domain"
	"github.com/gabrielaraujr/golang-case/account/internal/ports"
	"github.com/google/uuid"
)

type GetProposalUseCase struct {
	repository ports.ProposalRepository
}

func NewGetProposalUseCase(repo ports.ProposalRepository) *GetProposalUseCase {
	return &GetProposalUseCase{repository: repo}
}

func (uc *GetProposalUseCase) Execute(ctx context.Context, id uuid.UUID) (*dto.ProposalResponse, error) {
	proposal, err := uc.repository.FindByID(ctx, id)
	if err != nil && errors.Is(err, domainErrors.ErrProposalNotFound) {
		return nil, appErrors.NewNotFoundError("proposal")
	}
	if err != nil {
		return nil, appErrors.NewInternalError("failed to fetch proposal", err)
	}

	return entityToResponse(proposal), nil
}
