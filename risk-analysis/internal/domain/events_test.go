package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestProposalCreatedEventValidate(t *testing.T) {
	tests := []struct {
		name    string
		event   *ProposalCreatedEvent
		wantErr error
	}{
		{
			name: "valid event",
			event: &ProposalCreatedEvent{
				EventType:  EventProposalCreated,
				ProposalID: uuid.New(),
				Payload: &ProposalPayload{
					FullName: "John Doe",
					CPF:      "12345678902",
					Salary:   5000.0,
				},
			},
			wantErr: nil,
		},
		{
			name: "nil payload",
			event: &ProposalCreatedEvent{
				EventType:  EventProposalCreated,
				ProposalID: uuid.New(),
				Payload:    nil,
			},
			wantErr: ErrNilPayload,
		},
		{
			name: "empty full name",
			event: &ProposalCreatedEvent{
				EventType:  EventProposalCreated,
				ProposalID: uuid.New(),
				Payload: &ProposalPayload{
					FullName: "",
					CPF:      "12345678902",
					Salary:   5000.0,
				},
			},
			wantErr: ErrEmptyFullName,
		},
		{
			name: "empty cpf",
			event: &ProposalCreatedEvent{
				EventType:  EventProposalCreated,
				ProposalID: uuid.New(),
				Payload: &ProposalPayload{
					FullName: "John Doe",
					CPF:      "",
					Salary:   5000.0,
				},
			},
			wantErr: ErrEmptyCPF,
		},
		{
			name: "negative salary",
			event: &ProposalCreatedEvent{
				EventType:  EventProposalCreated,
				ProposalID: uuid.New(),
				Payload: &ProposalPayload{
					FullName: "John Doe",
					CPF:      "12345678902",
					Salary:   -100.0,
				},
			},
			wantErr: ErrNegativeSalary,
		},
		{
			name: "zero salary is valid",
			event: &ProposalCreatedEvent{
				EventType:  EventProposalCreated,
				ProposalID: uuid.New(),
				Payload: &ProposalPayload{
					FullName: "John Doe",
					CPF:      "12345678902",
					Salary:   0.0,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if err != tt.wantErr {
				t.Errorf("ProposalCreatedEvent.Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
