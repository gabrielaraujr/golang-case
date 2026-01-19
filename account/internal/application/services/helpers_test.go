package services

import (
	"context"
	"errors"
	"testing"

	appErrors "github.com/gabrielaraujr/golang-case/account/internal/application"
	"github.com/gabrielaraujr/golang-case/account/internal/application/dto"
	"github.com/gabrielaraujr/golang-case/account/internal/domain/entities"
	"github.com/gabrielaraujr/golang-case/account/internal/ports"
	"github.com/google/uuid"
)

type mockRepository struct {
	saveFn      func(ctx context.Context, p *entities.Proposal) error
	findByCPFFn func(ctx context.Context, cpf string) (*entities.Proposal, error)
	findByIDFn  func(ctx context.Context, id uuid.UUID) (*entities.Proposal, error)
}

func (m *mockRepository) Save(ctx context.Context, p *entities.Proposal) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, p)
	}
	return nil
}

func (m *mockRepository) FindByCPF(ctx context.Context, cpf string) (*entities.Proposal, error) {
	if m.findByCPFFn != nil {
		return m.findByCPFFn(ctx, cpf)
	}
	return nil, nil
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Proposal, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) Update(ctx context.Context, p *entities.Proposal) error {
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
	salary    float64
	email     string
	phone     string
	birthDate string
	address   dto.AddressRequest
}

func newRequestBuilder() *requestBuilder {
	return &requestBuilder{
		fullName:  "John Doe",
		cpf:       "12345678901",
		salary:    5000.00,
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
		Salary:    b.salary,
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
