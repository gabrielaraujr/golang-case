package handler

import (
	"context"

	"github.com/gabrielaraujr/golang-case/account/internal/application/dto"
	"github.com/google/uuid"
)

type mockCreateProposalUseCase struct {
	executeFn func(ctx context.Context, req *dto.CreateProposalRequest) (*dto.ProposalResponse, error)
}

func (m *mockCreateProposalUseCase) Execute(ctx context.Context, req *dto.CreateProposalRequest) (*dto.ProposalResponse, error) {
	if m.executeFn != nil {
		return m.executeFn(ctx, req)
	}
	return nil, nil
}

type mockGetProposalUseCase struct {
	executeFn func(ctx context.Context, id uuid.UUID) (*dto.ProposalResponse, error)
}

func (m *mockGetProposalUseCase) Execute(ctx context.Context, id uuid.UUID) (*dto.ProposalResponse, error) {
	if m.executeFn != nil {
		return m.executeFn(ctx, id)
	}
	return nil, nil
}
