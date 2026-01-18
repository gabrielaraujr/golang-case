package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

type ProposalBuilder struct {
	fullName  string
	cpf       string
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
	return NewProposal(b.fullName, b.cpf, b.email, b.phone, b.birthDate, b.address)
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

func assertErrorMessage(t *testing.T, err error, expected string) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error %q, got nil", expected)
		return
	}
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
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
		name    string
		builder *ProposalBuilder
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid proposal",
			builder: NewProposalBuilder(),
			wantErr: false,
		},
		{
			name:    "empty full name",
			builder: NewProposalBuilder().WithFullName(""),
			wantErr: true,
			errMsg:  "full name is required",
		},
		{
			name:    "empty CPF",
			builder: NewProposalBuilder().WithCPF(""),
			wantErr: true,
			errMsg:  "CPF is required",
		},
		{
			name:    "empty email",
			builder: NewProposalBuilder().WithEmail(""),
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name:    "empty phone",
			builder: NewProposalBuilder().WithPhone(""),
			wantErr: false,
		},
		{
			name:    "zero birth date",
			builder: NewProposalBuilder().WithBirthDate(time.Time{}),
			wantErr: true,
			errMsg:  "birth date is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposal, err := tt.builder.BuildWithValidation()

			if tt.wantErr {
				assertErrorMessage(t, err, tt.errMsg)
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
	t.Run("start analysis - pending to analyzing", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		assertNoError(t, p.StartAnalysis())
		assertStatus(t, p.Status, StatusAnalyzing)
	})

	t.Run("start analysis - invalid from analyzing", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusAnalyzing).Build()
		assertError(t, p.StartAnalysis())
	})

	t.Run("approve - analyzing to approved", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusAnalyzing).Build()
		assertNoError(t, p.Approve())
		assertStatus(t, p.Status, StatusApproved)
	})

	t.Run("approve - invalid from pending", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		assertError(t, p.Approve())
	})

	t.Run("reject - pending to rejected", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		assertNoError(t, p.Reject())
		assertStatus(t, p.Status, StatusRejected)
	})

	t.Run("reject - analyzing to rejected", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusAnalyzing).Build()
		assertNoError(t, p.Reject())
		assertStatus(t, p.Status, StatusRejected)
	})

	t.Run("reject - invalid from approved", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusApproved).Build()
		assertError(t, p.Reject())
	})
}

func TestProposalStatusQueries(t *testing.T) {
	t.Run("IsPending", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		assertBool(t, p.IsPending(), "expected IsPending to return true")
		assertBool(t, !p.IsAnalyzing() && !p.IsFinalized(), "expected IsAnalyzing and IsFinalized to return false")
	})

	t.Run("IsAnalyzing", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusAnalyzing).Build()
		assertBool(t, p.IsAnalyzing(), "expected IsAnalyzing to return true")
		assertBool(t, !p.IsPending() && !p.IsFinalized(), "expected IsPending and IsFinalized to return false")
	})

	t.Run("IsFinalized - approved", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusApproved).Build()
		assertBool(t, p.IsFinalized(), "expected IsFinalized to return true for approved status")
	})

	t.Run("IsFinalized - rejected", func(t *testing.T) {
		p := NewProposalBuilder().WithStatus(StatusRejected).Build()
		assertBool(t, p.IsFinalized(), "expected IsFinalized to return true for rejected status")
	})
}

func TestProposalIsValid(t *testing.T) {
	t.Run("valid proposal", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		assertBool(t, p.IsValid(), "expected IsValid to return true for valid proposal")
	})

	t.Run("invalid - nil UUID", func(t *testing.T) {
		p := NewProposalBuilder().Build()
		p.ID = uuid.Nil
		assertBool(t, !p.IsValid(), "expected IsValid to return false for nil UUID")
	})

	t.Run("invalid - empty fields", func(t *testing.T) {
		tests := []struct {
			name    string
			builder *ProposalBuilder
		}{
			{"empty full name", NewProposalBuilder().WithFullName("")},
			{"empty CPF", NewProposalBuilder().WithCPF("")},
			{"empty email", NewProposalBuilder().WithEmail("")},
			{"empty phone", NewProposalBuilder().WithPhone("")},
			{"zero birth date", NewProposalBuilder().WithBirthDate(time.Time{})},
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
