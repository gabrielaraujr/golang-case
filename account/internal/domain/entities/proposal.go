package entities

import (
	"time"

	errors "github.com/gabrielaraujr/golang-case/account/internal/domain"
	"github.com/google/uuid"
)

type ProposalStatus string

const (
	StatusPending   ProposalStatus = "pending"
	StatusAnalyzing ProposalStatus = "analyzing"
	StatusApproved  ProposalStatus = "approved"
	StatusRejected  ProposalStatus = "rejected"
)

type Proposal struct {
	ID        uuid.UUID
	FullName  string
	CPF       string
	Salary    float64
	BirthDate time.Time
	Email     string
	Phone     string
	Address   Address
	Status    ProposalStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

func NewProposal(fullName, cpf string, salary float64, email, phone string, birthDate time.Time, address Address) (*Proposal, error) {
	if fullName == "" {
		return nil, errors.ErrFullNameRequired
	}
	if cpf == "" {
		return nil, errors.ErrCPFRequired
	}
	if salary <= 0 {
		return nil, errors.ErrSalaryRequired
	}
	if email == "" {
		return nil, errors.ErrEmailRequired
	}
	if birthDate.IsZero() {
		return nil, errors.ErrBirthDateRequired
	}

	return &Proposal{
		ID:        uuid.New(),
		FullName:  fullName,
		CPF:       cpf,
		Salary:    salary,
		Email:     email,
		Phone:     phone,
		BirthDate: birthDate,
		Address:   address,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (p *Proposal) Approve() error {
	if p.Status != StatusAnalyzing {
		return errors.ErrOnlyAnalyzingCanBeApproved
	}
	p.Status = StatusApproved
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Proposal) StartAnalysis() error {
	if p.Status != StatusPending {
		return errors.ErrOnlyPendingCanStartAnalysis
	}
	p.Status = StatusAnalyzing
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Proposal) Reject() error {
	if p.Status != StatusPending && p.Status != StatusAnalyzing {
		return errors.ErrOnlyPendingOrAnalyzingCanReject
	}
	p.Status = StatusRejected
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Proposal) IsPending() bool {
	return p.Status == StatusPending
}

func (p *Proposal) IsAnalyzing() bool {
	return p.Status == StatusAnalyzing
}

func (p *Proposal) IsFinalized() bool {
	return p.Status == StatusApproved || p.Status == StatusRejected
}

func (p *Proposal) IsValid() bool {
	return p.ID != uuid.Nil &&
		p.FullName != "" &&
		p.CPF != "" &&
		p.Salary != 0 &&
		p.Email != "" &&
		p.Phone != "" &&
		!p.BirthDate.IsZero() &&
		p.Address.Street != "" &&
		p.Address.City != "" &&
		p.Address.State != "" &&
		p.Address.ZipCode != "" &&
		p.Status != ""
}
