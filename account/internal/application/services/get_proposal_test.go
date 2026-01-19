package services

import (
	"context"
	"errors"
	"testing"
	"time"

	domainErrors "github.com/gabrielaraujr/golang-case/account/internal/domain"
	"github.com/gabrielaraujr/golang-case/account/internal/domain/entities"
	"github.com/google/uuid"
)

func TestGetProposalUseCase_Execute(t *testing.T) {
	t.Run("should return proposal when found", func(t *testing.T) {
		expectedProposal := &entities.Proposal{
			ID:        uuid.New(),
			FullName:  "John Doe",
			CPF:       "12345678901",
			Email:     "john@example.com",
			Phone:     "11999999999",
			BirthDate: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			Address: entities.Address{
				Street:  "123 Main St",
				City:    "São Paulo",
				State:   "SP",
				ZipCode: "01234-567",
			},
			Status:    entities.StatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		repo := &mockRepository{
			findByIDFn: func(ctx context.Context, id uuid.UUID) (*entities.Proposal, error) {
				return expectedProposal, nil
			},
		}

		useCase := NewGetProposalUseCase(repo)
		response, err := useCase.Execute(context.Background(), expectedProposal.ID)

		assertNoError(t, err)
		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.ID != expectedProposal.ID {
			t.Errorf("expected ID %v, got %v", expectedProposal.ID, response.ID)
		}
		if response.CPF != expectedProposal.CPF {
			t.Errorf("expected CPF %q, got %q", expectedProposal.CPF, response.CPF)
		}
		if response.Status != string(expectedProposal.Status) {
			t.Errorf("expected status %q, got %q", expectedProposal.Status, response.Status)
		}
	})

	t.Run("should return not found error when proposal does not exist", func(t *testing.T) {
		repo := &mockRepository{
			findByIDFn: func(ctx context.Context, id uuid.UUID) (*entities.Proposal, error) {
				return nil, domainErrors.ErrProposalNotFound
			},
		}

		useCase := NewGetProposalUseCase(repo)
		proposalID := uuid.New()

		response, err := useCase.Execute(context.Background(), proposalID)

		assertError(t, err)
		assertApplicationError(t, err, "NOT_FOUND", 404)
		if response != nil {
			t.Error("expected nil response")
		}
	})

	t.Run("should return internal error when repository fails", func(t *testing.T) {
		repo := &mockRepository{
			findByIDFn: func(ctx context.Context, id uuid.UUID) (*entities.Proposal, error) {
				return nil, errors.New("database connection error")
			},
		}

		useCase := NewGetProposalUseCase(repo)
		proposalID := uuid.New()

		response, err := useCase.Execute(context.Background(), proposalID)

		assertError(t, err)
		assertApplicationError(t, err, "INTERNAL_ERROR", 500)
		if response != nil {
			t.Error("expected nil response")
		}
	})
}
