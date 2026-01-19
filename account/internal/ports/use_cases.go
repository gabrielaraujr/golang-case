package ports

import (
	"context"

	"github.com/gabrielaraujr/golang-case/account/internal/application/dto"
	"github.com/google/uuid"
)

type CreateProposalUseCase interface {
	Execute(ctx context.Context, req *dto.CreateProposalRequest) (*dto.ProposalResponse, error)
}

type GetProposalUseCase interface {
	Execute(ctx context.Context, id uuid.UUID) (*dto.ProposalResponse, error)
}
