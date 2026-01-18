package services

import (
	"context"
	"errors"
	"testing"
	"time"

	appErrors "github.com/gabrielaraujr/golang-case/account/internal/application"
	"github.com/gabrielaraujr/golang-case/account/internal/application/dto"
	"github.com/gabrielaraujr/golang-case/account/internal/domain/entities"
	"github.com/gabrielaraujr/golang-case/account/internal/ports"
	"github.com/google/uuid"
)

type mockProposalRepository struct {
	saveFn      func(ctx context.Context, p *entities.Proposal) error
	findByCPFFn func(ctx context.Context, cpf string) (*entities.Proposal, error)
}

func (m *mockProposalRepository) Save(ctx context.Context, p *entities.Proposal) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, p)
	}
	return nil
}

func (m *mockProposalRepository) FindByCPF(ctx context.Context, cpf string) (*entities.Proposal, error) {
	if m.findByCPFFn != nil {
		return m.findByCPFFn(ctx, cpf)
	}
	return nil, nil
}

func (m *mockProposalRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Proposal, error) {
	return nil, nil
}

func (m *mockProposalRepository) Update(ctx context.Context, p *entities.Proposal) error {
	return nil
}

type mockQueueProducer struct {
	publishFn func(ctx context.Context, event *ports.ProposalEvent) error
}

func (m *mockQueueProducer) Publish(ctx context.Context, event *ports.ProposalEvent) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, event)
	}
	return nil
}

type mockLogger struct {
	infoFn  func(ctx context.Context, msg string, args ...interface{})
	errorFn func(ctx context.Context, msg string, args ...interface{})
}

func (m *mockLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if m.infoFn != nil {
		m.infoFn(ctx, msg, args...)
	}
}

func (m *mockLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if m.errorFn != nil {
		m.errorFn(ctx, msg, args...)
	}
}

func (m *mockLogger) Warn(ctx context.Context, msg string, args ...interface{}) {}

type requestBuilder struct {
	fullName  string
	cpf       string
	email     string
	phone     string
	birthDate string
	address   dto.AddressRequest
}

func newRequestBuilder() *requestBuilder {
	return &requestBuilder{
		fullName:  "John Doe",
		cpf:       "12345678901",
		email:     "john@example.com",
		phone:     "11999999999",
		birthDate: "15-01-1990",
		address: dto.AddressRequest{
			Street:  "123 Main St",
			City:    "São Paulo",
			State:   "SP",
			ZipCode: "01234-567",
		},
	}
}

func (b *requestBuilder) withCPF(cpf string) *requestBuilder {
	b.cpf = cpf
	return b
}

func (b *requestBuilder) withBirthDate(date string) *requestBuilder {
	b.birthDate = date
	return b
}

func (b *requestBuilder) build() *dto.CreateProposalRequest {
	return &dto.CreateProposalRequest{
		FullName:  b.fullName,
		CPF:       b.cpf,
		Email:     b.email,
		Phone:     b.phone,
		BirthDate: b.birthDate,
		Address:   b.address,
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func assertApplicationError(t *testing.T, err error, expectedCode string, expectedStatus int) {
	t.Helper()
	var appErr *appErrors.ApplicationError
	if !errors.As(err, &appErr) {
		t.Errorf("expected ApplicationError, got %T", err)
		return
	}
	if appErr.Code != expectedCode {
		t.Errorf("expected code %q, got %q", expectedCode, appErr.Code)
	}
	if appErr.StatusCode != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, appErr.StatusCode)
	}
}

func TestCreateProposalUseCase_Execute(t *testing.T) {
	t.Run("should create proposal successfully", func(t *testing.T) {
		repo := &mockProposalRepository{}
		producer := &mockQueueProducer{}
		logger := &mockLogger{}

		useCase := NewCreateProposalUseCase(repo, producer, logger)
		req := newRequestBuilder().build()

		response, err := useCase.Execute(context.Background(), req)

		assertNoError(t, err)
		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.ID == uuid.Nil {
			t.Error("expected valid UUID")
		}
		if response.CPF != req.CPF {
			t.Errorf("expected CPF %q, got %q", req.CPF, response.CPF)
		}
		if response.Status != string(entities.StatusPending) {
			t.Errorf("expected status %q, got %q", entities.StatusPending, response.Status)
		}
	})

	t.Run("should return error for invalid birth date format", func(t *testing.T) {
		repo := &mockProposalRepository{}
		producer := &mockQueueProducer{}
		logger := &mockLogger{}

		useCase := NewCreateProposalUseCase(repo, producer, logger)
		req := newRequestBuilder().withBirthDate("1990-01-15").build()

		response, err := useCase.Execute(context.Background(), req)

		assertError(t, err)
		assertApplicationError(t, err, "INVALID_INPUT", 400)
		if response != nil {
			t.Error("expected nil response")
		}
	})

	t.Run("should return error when CPF already exists", func(t *testing.T) {
		existingProposal := &entities.Proposal{
			ID:  uuid.New(),
			CPF: "12345678901",
		}

		repo := &mockProposalRepository{
			findByCPFFn: func(ctx context.Context, cpf string) (*entities.Proposal, error) {
				return existingProposal, nil
			},
		}
		producer := &mockQueueProducer{}
		logger := &mockLogger{}

		useCase := NewCreateProposalUseCase(repo, producer, logger)
		req := newRequestBuilder().withCPF("12345678901").build()

		response, err := useCase.Execute(context.Background(), req)

		assertError(t, err)
		assertApplicationError(t, err, "DUPLICATE_CPF", 409)
		if response != nil {
			t.Error("expected nil response")
		}
	})

	t.Run("should return error when repository save fails", func(t *testing.T) {
		repo := &mockProposalRepository{
			saveFn: func(ctx context.Context, p *entities.Proposal) error {
				return errors.New("database error")
			},
		}
		producer := &mockQueueProducer{}
		logger := &mockLogger{}

		useCase := NewCreateProposalUseCase(repo, producer, logger)
		req := newRequestBuilder().build()

		response, err := useCase.Execute(context.Background(), req)

		assertError(t, err)
		assertApplicationError(t, err, "INTERNAL_ERROR", 500)
		if response != nil {
			t.Error("expected nil response")
		}
	})

	t.Run("should publish event after successful creation", func(t *testing.T) {
		var publishedEvent *ports.ProposalEvent

		repo := &mockProposalRepository{}
		producer := &mockQueueProducer{
			publishFn: func(ctx context.Context, event *ports.ProposalEvent) error {
				publishedEvent = event
				return nil
			},
		}
		logger := &mockLogger{}

		useCase := NewCreateProposalUseCase(repo, producer, logger)
		req := newRequestBuilder().build()

		response, err := useCase.Execute(context.Background(), req)

		assertNoError(t, err)
		if publishedEvent == nil {
			t.Fatal("expected event to be published")
		}
		if publishedEvent.EventType != "ProposalCreated" {
			t.Errorf("expected event type %q, got %q", "ProposalCreated", publishedEvent.EventType)
		}
		if publishedEvent.ProposalID != response.ID.String() {
			t.Error("event proposal ID doesn't match response ID")
		}
	})

	t.Run("should continue when event publishing fails", func(t *testing.T) {
		repo := &mockProposalRepository{}
		producer := &mockQueueProducer{
			publishFn: func(ctx context.Context, event *ports.ProposalEvent) error {
				return errors.New("queue error")
			},
		}
		logger := &mockLogger{}

		useCase := NewCreateProposalUseCase(repo, producer, logger)
		req := newRequestBuilder().build()

		response, err := useCase.Execute(context.Background(), req)

		assertNoError(t, err)
		if response == nil {
			t.Error("expected response even when event publishing fails")
		}
	})
}

func TestEntityToResponse(t *testing.T) {
	proposal := &entities.Proposal{
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

	response := entityToResponse(proposal)

	if response.ID != proposal.ID {
		t.Error("ID mismatch")
	}
	if response.CPF != proposal.CPF {
		t.Error("CPF mismatch")
	}
	if response.Status != string(proposal.Status) {
		t.Error("Status mismatch")
	}
	if response.Address.City != proposal.Address.City {
		t.Error("Address city mismatch")
	}
}
