package entities

import (
	"errors"
	"time"

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

func NewProposal(fullName, cpf, email, phone string, birthDate time.Time, address Address) (*Proposal, error) {
	if fullName == "" {
		return nil, errors.New("full name is required")
	}
	if cpf == "" {
		return nil, errors.New("CPF is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if birthDate.IsZero() {
		return nil, errors.New("birth date is required")
	}

	return &Proposal{
		ID:        uuid.New(),
		FullName:  fullName,
		CPF:       cpf,
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
		return errors.New("only analyzing proposals can be approved")
	}
	p.Status = StatusApproved
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Proposal) StartAnalysis() error {
	if p.Status != StatusPending {
		return errors.New("only pending proposals can start analysis")
	}
	p.Status = StatusAnalyzing
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Proposal) Reject() error {
	if p.Status != StatusPending && p.Status != StatusAnalyzing {
		return errors.New("only pending or analyzing proposals can be rejected")
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
		p.Email != "" &&
		p.Phone != "" &&
		!p.BirthDate.IsZero() &&
		p.Address.Street != "" &&
		p.Address.City != "" &&
		p.Address.State != "" &&
		p.Address.ZipCode != "" &&
		p.Status != ""
}
