package services

import (
	"context"
	"time"

	appErrors "github.com/gabrielaraujr/golang-case/account/internal/application"
	"github.com/gabrielaraujr/golang-case/account/internal/application/dto"
	"github.com/gabrielaraujr/golang-case/account/internal/domain/entities"
	"github.com/gabrielaraujr/golang-case/account/internal/ports"
)

type CreateProposalUseCase struct {
	repository ports.ProposalRepository
	producer   ports.QueueProducer
	logger     ports.Logger
}

const DateLayoutBR = "02-01-2006" // Brazilian format (dd-mm-yyyy)

func NewCreateProposalUseCase(
	repo ports.ProposalRepository,
	prod ports.QueueProducer,
	logger ports.Logger,
) *CreateProposalUseCase {
	return &CreateProposalUseCase{
		repository: repo,
		producer:   prod,
		logger:     logger,
	}
}

func (uc *CreateProposalUseCase) Execute(
	ctx context.Context,
	req *dto.CreateProposalRequest,
) (*dto.ProposalResponse, error) {
	uc.logger.Info(ctx, "creating proposal", "cpf", req.CPF)

	birthDate, err := time.Parse(DateLayoutBR, req.BirthDate)
	if err != nil {
		return nil, appErrors.NewInvalidInputError(err)
	}

	address := entities.Address{
		Street:  req.Address.Street,
		City:    req.Address.City,
		State:   req.Address.State,
		ZipCode: req.Address.ZipCode,
	}

	proposal, err := entities.NewProposal(
		req.FullName,
		req.CPF,
		req.Email,
		req.Phone,
		birthDate,
		address,
	)
	if err != nil {
		return nil, appErrors.NewInvalidInputError(err)
	}

	existing, _ := uc.repository.FindByCPF(ctx, req.CPF)
	if existing != nil {
		return nil, appErrors.NewDuplicateCPFError()
	}

	if err := uc.repository.Save(ctx, proposal); err != nil {
		uc.logger.Error(ctx, "failed to save proposal", "error", err)
		return nil, appErrors.NewInternalError("failed to save proposal", err)
	}

	event := &ports.ProposalEvent{
		EventType:  "ProposalCreated",
		ProposalID: proposal.ID.String(),
		Payload: map[string]interface{}{
			"full_name": proposal.FullName,
			"cpf":       proposal.CPF,
			"email":     proposal.Email,
		},
	}
	_ = uc.producer.Publish(ctx, event) // Fire and forget

	uc.logger.Info(ctx, "proposal created", "proposal_id", proposal.ID.String())
	return entityToResponse(proposal), nil
}

func entityToResponse(p *entities.Proposal) *dto.ProposalResponse {
	return &dto.ProposalResponse{
		ID:        p.ID,
		FullName:  p.FullName,
		CPF:       p.CPF,
		Email:     p.Email,
		Phone:     p.Phone,
		BirthDate: p.BirthDate,
		Address: dto.AddressResponse{
			Street:  p.Address.Street,
			City:    p.Address.City,
			State:   p.Address.State,
			ZipCode: p.Address.ZipCode,
		},
		Status:    string(p.Status),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
