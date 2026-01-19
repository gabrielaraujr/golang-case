package entities

import (
	"errors"
	"testing"
	"time"

	domainErrors "github.com/gabrielaraujr/golang-case/account/internal/domain"
	"github.com/google/uuid"
)

type ProposalBuilder struct {
	fullName  string
	cpf       string
	salary    float64
	email     string
	phone     string
	birthDate time.Time
	address   Address
	status    ProposalStatus
}

func NewProposalBuilder() *ProposalBuilder {
	return &ProposalBuilder{
		fullName:  "John Doe",
		cpf:       "12345678901",
		salary:    5000.00,
		email:     "john@example.com",
		phone:     "11999999999",
		birthDate: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		address: Address{
			Street:  "123 Main St",
			City:    "São Paulo",
			State:   "SP",
			ZipCode: "01234-567",
		},
		status: StatusPending,
	}
}

func (b *ProposalBuilder) WithFullName(name string) *ProposalBuilder {
	b.fullName = name
	return b
}

func (b *ProposalBuilder) WithCPF(cpf string) *ProposalBuilder {
	b.cpf = cpf
	return b
}

func (b *ProposalBuilder) WithSalary(salary float64) *ProposalBuilder {
	b.salary = salary
	return b
}

func (b *ProposalBuilder) WithEmail(email string) *ProposalBuilder {
	b.email = email
	return b
}

func (b *ProposalBuilder) WithPhone(phone string) *ProposalBuilder {
	b.phone = phone
	return b
}

func (b *ProposalBuilder) WithBirthDate(birthDate time.Time) *ProposalBuilder {
	b.birthDate = birthDate
	return b
}

func (b *ProposalBuilder) WithAddress(address Address) *ProposalBuilder {
	b.address = address
	return b
}

func (b *ProposalBuilder) WithStatus(status ProposalStatus) *ProposalBuilder {
	b.status = status
	return b
}

func (b *ProposalBuilder) Build() *Proposal {
	return &Proposal{
		ID:        uuid.New(),
		FullName:  b.fullName,
		CPF:       b.cpf,
		Salary:    b.salary,
		Email:     b.email,
		Phone:     b.phone,
		BirthDate: b.birthDate,
		Address:   b.address,
		Status:    b.status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (b *ProposalBuilder) BuildWithValidation() (*Proposal, error) {
	return NewProposal(b.fullName, b.cpf, b.salary, b.email, b.phone, b.birthDate, b.address)
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

func assertErrorIs(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("expected error %v, got %v", want, got)
	}
}

func assertStatus(t *testing.T, got, want ProposalStatus) {
	t.Helper()
	if got != want {
		t.Errorf("expected status %q, got %q", want, got)
	}
}

func assertBool(t *testing.T, got bool, message string) {
	t.Helper()
	if !got {
		t.Error(message)
	}
}

func TestNewProposal(t *testing.T) {
	tests := []struct {
		name        string
		builder     *ProposalBuilder
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "should create proposal with valid data",
			builder: NewProposalBuilder(),
			wantErr: false,
		},
		{
			name:        "should return error when full name is empty",
			builder:     NewProposalBuilder().WithFullName(""),
			wantErr:     true,
			expectedErr: domainErrors.ErrFullNameRequired,
		},
		{
			name:        "should return error when CPF is empty",
			builder:     NewProposalBuilder().WithCPF(""),
			wantErr:     true,
			expectedErr: domainErrors.ErrCPFRequired,
		},
		{
			name:        "should return error when salary is zero",
			builder:     NewProposalBuilder().WithSalary(0),
			wantErr:     true,
			expectedErr: domainErrors.ErrSalaryRequired,
		},
		{
			name:        "should return error when email is empty",
			builder:     NewProposalBuilder().WithEmail(""),
			wantErr:     true,
			expectedErr: domainErrors.ErrEmailRequired,
		},
		{
			name:    "should allow empty phone number",
			builder: NewProposalBuilder().WithPhone(""),
			wantErr: false,
		},
		{
			name:        "should return error when birth date is zero",
			builder:     NewProposalBuilder().WithBirthDate(time.Time{}),
			wantErr:     true,
			expectedErr: domainErrors.ErrBirthDateRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposal, err := tt.builder.BuildWithValidation()

			if tt.wantErr {
				assertError(t, err)
				assertErrorIs(t, err, tt.expectedErr)
				return
			}

			assertNoError(t, err)
			if proposal.ID == uuid.Nil {
				t.Error("expected non-nil UUID")
			}
			assertStatus(t, proposal.Status, StatusPending)
		})
	}
}

func TestProposalStateTransitions(t *testing.T) {
	t.Run("should transition from pending to analyzing", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		assertNoError(t, p.StartAnalysis())
		assertStatus(t, p.Status, StatusAnalyzing)
	})

	t.Run("should return error when starting analysis from analyzing status", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusAnalyzing).Build()
		err := p.StartAnalysis()
		assertError(t, err)
		assertErrorIs(t, err, domainErrors.ErrOnlyPendingCanStartAnalysis)
	})

	t.Run("should return error when starting analysis from approved status", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusApproved).Build()
		err := p.StartAnalysis()
		assertError(t, err)
		assertErrorIs(t, err, domainErrors.ErrOnlyPendingCanStartAnalysis)
	})

	t.Run("should transition from analyzing to approved", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusAnalyzing).Build()
		assertNoError(t, p.Approve())
		assertStatus(t, p.Status, StatusApproved)
	})

	t.Run("should return error when approving from pending status", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		err := p.Approve()
		assertError(t, err)
		assertErrorIs(t, err, domainErrors.ErrOnlyAnalyzingCanBeApproved)
	})

	t.Run("should return error when approving from rejected status", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusRejected).Build()
		err := p.Approve()
		assertError(t, err)
		assertErrorIs(t, err, domainErrors.ErrOnlyAnalyzingCanBeApproved)
	})

	t.Run("should transition from pending to rejected", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		assertNoError(t, p.Reject())
		assertStatus(t, p.Status, StatusRejected)
	})

	t.Run("should transition from analyzing to rejected", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusAnalyzing).Build()
		assertNoError(t, p.Reject())
		assertStatus(t, p.Status, StatusRejected)
	})

	t.Run("should return error when rejecting from approved status", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusApproved).Build()
		err := p.Reject()
		assertError(t, err)
		assertErrorIs(t, err, domainErrors.ErrOnlyPendingOrAnalyzingCanReject)
	})

	t.Run("should return error when rejecting from rejected status", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusRejected).Build()
		err := p.Reject()
		assertError(t, err)
		assertErrorIs(t, err, domainErrors.ErrOnlyPendingOrAnalyzingCanReject)
	})
}

func TestProposalStatusQueries(t *testing.T) {
	t.Run("should return true for IsPending when status is pending", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		assertBool(t, p.IsPending(), "expected IsPending to return true")
		assertBool(t, !p.IsAnalyzing() && !p.IsFinalized(), "expected IsAnalyzing and IsFinalized to return false")
	})

	t.Run("should return true for IsAnalyzing when status is analyzing", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusAnalyzing).Build()
		assertBool(t, p.IsAnalyzing(), "expected IsAnalyzing to return true")
		assertBool(t, !p.IsPending() && !p.IsFinalized(), "expected IsPending and IsFinalized to return false")
	})

	t.Run("should return true for IsFinalized when status is approved", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusApproved).Build()
		assertBool(t, p.IsFinalized(), "expected IsFinalized to return true for approved status")
	})

	t.Run("should return true for IsFinalized when status is rejected", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusRejected).Build()
		assertBool(t, p.IsFinalized(), "expected IsFinalized to return true for rejected status")
	})
}

func TestProposalIsValid(t *testing.T) {
	t.Run("should return true when proposal has all required fields", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		assertBool(t, p.IsValid(), "expected IsValid to return true for valid proposal")
	})

	t.Run("should return false when UUID is nil", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		p.ID = uuid.Nil
		assertBool(t, !p.IsValid(), "expected IsValid to return false for nil UUID")
	})

	t.Run("should return false when required fields are empty", func(t *testing.T) {
		tests := []struct {
			name    string
			builder *ProposalBuilder
		}{
			{"empty full name", NewProposalBuilder().WithFullName("")},
			{"empty CPF", NewProposalBuilder().WithCPF("")},
			{"empty salary", NewProposalBuilder().WithSalary(0)},
			{"empty email", NewProposalBuilder().WithEmail("")},
			{"empty phone", NewProposalBuilder().WithPhone("")},
			{"empty birth date", NewProposalBuilder().WithBirthDate(time.Time{})},
			{"empty address street", NewProposalBuilder().WithAddress(Address{City: "SP", State: "SP", ZipCode: "01234-567"})},
			{"empty address city", NewProposalBuilder().WithAddress(Address{Street: "123 Main", State: "SP", ZipCode: "01234-567"})},
			{"empty address state", NewProposalBuilder().WithAddress(Address{Street: "123 Main", City: "SP", ZipCode: "01234-567"})},
			{"empty address zip", NewProposalBuilder().WithAddress(Address{Street: "123 Main", City: "SP", State: "SP"})},
			{"empty status", NewProposalBuilder().WithStatus("")},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				p := tt.builder.Build()
				assertBool(t, !p.IsValid(), "expected IsValid to return false")
			})
		}
	})
}
