package services

import (
	"context"
	"testing"

	events "github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain"
	"github.com/google/uuid"
)

type mockQueueProducer struct {
	publishFunc func(ctx context.Context, event *events.ProposalStatusChangedEvent) error
	published   []*events.ProposalStatusChangedEvent
}

func newMockQueueProducer() *mockQueueProducer {
	return &mockQueueProducer{
		published: make([]*events.ProposalStatusChangedEvent, 0),
	}
}

func (m *mockQueueProducer) Publish(ctx context.Context, event *events.ProposalStatusChangedEvent) error {
	m.published = append(m.published, event)
	if m.publishFunc != nil {
		return m.publishFunc(ctx, event)
	}
	return nil
}

type mockLogger struct {
	infoCalls  int
	errorCalls int
	warnCalls  int
}

func newMockLogger() *mockLogger {
	return &mockLogger{}
}

func (m *mockLogger) Info(ctx context.Context, msg string, args ...any) {
	m.infoCalls++
}

func (m *mockLogger) Error(ctx context.Context, msg string, args ...any) {
	m.errorCalls++
}

func (m *mockLogger) Warn(ctx context.Context, msg string, args ...any) {
	m.warnCalls++
}

func assertEventCount(t *testing.T, events []*events.ProposalStatusChangedEvent, expected int) {
	t.Helper()
	if len(events) != expected {
		t.Fatalf("expected %d events, got %d", expected, len(events))
	}
}

func assertEvent(
	t *testing.T,
	event *events.ProposalStatusChangedEvent,
	expectedType string,
	expectedApproved bool,
	expectedProposalID uuid.UUID,
) {
	t.Helper()
	if event.EventType != expectedType {
		t.Errorf("event.EventType = %q, want %q", event.EventType, expectedType)
	}
	if event.Approved != expectedApproved {
		t.Errorf("event.Approved = %v, want %v", event.Approved, expectedApproved)
	}
	if event.ProposalID != expectedProposalID {
		t.Errorf("event.ProposalID = %v, want %v", event.ProposalID, expectedProposalID)
	}
}

func TestAnalyzeProposalServiceHandle(t *testing.T) {
	tests := []struct {
		name           string
		payload        *events.ProposalPayload
		wantEvents     int
		wantEventTypes []string
		wantApproved   []bool
	}{
		{
			name:           "documents rejection",
			payload:        &events.ProposalPayload{CPF: "123", FullName: "John Doe", Salary: 5000.0},
			wantEvents:     1,
			wantEventTypes: []string{events.EventDocumentsRejected},
			wantApproved:   []bool{false},
		},
		{
			name:           "credit rejection",
			payload:        &events.ProposalPayload{CPF: "12345678902", FullName: "John Doe", Salary: 2000.0},
			wantEvents:     2,
			wantEventTypes: []string{events.EventDocumentsApproved, events.EventCreditRejected},
			wantApproved:   []bool{true, false},
		},
		{
			name:           "fraud rejection",
			payload:        &events.ProposalPayload{CPF: "12345678901", FullName: "John Doe", Salary: 5000.0},
			wantEvents:     2,
			wantEventTypes: []string{events.EventDocumentsApproved, events.EventFraudRejected},
			wantApproved:   []bool{true, false},
		},
		{
			name:           "all approved",
			payload:        &events.ProposalPayload{CPF: "12345678902", FullName: "John Doe", Salary: 5000.0},
			wantEvents:     2,
			wantEventTypes: []string{events.EventDocumentsApproved, events.EventRiskAnalysisCompleted},
			wantApproved:   []bool{true, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queueProducer := newMockQueueProducer()
			logger := newMockLogger()
			service := NewAnalyzeProposalService(queueProducer, logger)
			ctx := context.Background()
			proposalID := uuid.New()

			event := &events.ProposalCreatedEvent{
				EventType:  events.EventProposalCreated,
				ProposalID: proposalID,
				Payload:    tt.payload,
			}

			_ = service.Handle(ctx, event)

			assertEventCount(t, queueProducer.published, tt.wantEvents)
			for i := 0; i < tt.wantEvents; i++ {
				assertEvent(t, queueProducer.published[i], tt.wantEventTypes[i], tt.wantApproved[i], proposalID)
			}
		})
	}
}

func TestAnalyzeProposalServiceHandleValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload *events.ProposalPayload
		wantErr error
	}{
		{
			name:    "nil payload",
			payload: nil,
			wantErr: events.ErrNilPayload,
		},
		{
			name: "empty full name",
			payload: &events.ProposalPayload{
				FullName: "",
				CPF:      "12345678902",
				Salary:   5000.0,
			},
			wantErr: events.ErrEmptyFullName,
		},
		{
			name: "empty cpf",
			payload: &events.ProposalPayload{
				FullName: "John Doe",
				CPF:      "",
				Salary:   5000.0,
			},
			wantErr: events.ErrEmptyCPF,
		},
		{
			name: "negative salary",
			payload: &events.ProposalPayload{
				FullName: "John Doe",
				CPF:      "12345678902",
				Salary:   -100.0,
			},
			wantErr: events.ErrNegativeSalary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queueProducer := newMockQueueProducer()
			logger := newMockLogger()
			service := NewAnalyzeProposalService(queueProducer, logger)
			ctx := context.Background()

			event := &events.ProposalCreatedEvent{
				EventType:  events.EventProposalCreated,
				ProposalID: uuid.New(),
				Payload:    tt.payload,
			}

			err := service.Handle(ctx, event)
			if err != tt.wantErr {
				t.Errorf("Handle() error = %v, want %v", err, tt.wantErr)
			}
			assertEventCount(t, queueProducer.published, 0)
		})
	}
}
