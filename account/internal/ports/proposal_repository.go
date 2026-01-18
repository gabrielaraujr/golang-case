package ports

import (
	"context"

	"github.com/gabrielaraujr/golang-case/account/internal/domain/entities"
	"github.com/google/uuid"
)

type ProposalRepository interface {
	Save(ctx context.Context, proposal *entities.Proposal) error
	Update(ctx context.Context, proposal *entities.Proposal) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Proposal, error)
	FindByCPF(ctx context.Context, cpf string) (*entities.Proposal, error)
}
