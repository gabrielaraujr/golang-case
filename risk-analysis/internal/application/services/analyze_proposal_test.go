package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain/entities"
	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain/events"
	"github.com/google/uuid"
)

type mockQueueProducer struct {
	publishFunc func(ctx context.Context, event *events.RiskEvent) error
	published   []*events.RiskEvent
}

func newMockQueueProducer() *mockQueueProducer {
	return &mockQueueProducer{
		published: make([]*events.RiskEvent, 0),
	}
}

func (m *mockQueueProducer) Publish(ctx context.Context, event *events.RiskEvent) error {
	m.published = append(m.published, event)
	if m.publishFunc != nil {
		return m.publishFunc(ctx, event)
	}
	return nil
}

func validPayload() *entities.ProposalPayload {
	return &entities.ProposalPayload{
		CPF:      "12345678902",
		FullName: "John Doe",
		Salary:   5000.0,
	}
}

func assertEventCount(t *testing.T, events []*events.RiskEvent, expected int) {
	t.Helper()
	if len(events) != expected {
		t.Fatalf("expected %d events, got %d", expected, len(events))
	}
}

func assertEvent(
	t *testing.T,
	event *events.RiskEvent,
	expectedType string,
	expectedApproved bool,
	expectedProposalID uuid.UUID,
) {
	t.Helper()
	if event.EventType != expectedType {
		t.Errorf("event.EventType = %q, expected %q", event.EventType, expectedType)
	}
	if event.Approved != expectedApproved {
		t.Errorf("event.Approved = %v, expected %v", event.Approved, expectedApproved)
	}
	if event.ProposalID != expectedProposalID {
		t.Errorf("event.ProposalID = %v, expected %v", event.ProposalID, expectedProposalID)
	}
}

func TestAnalyzeProposalServiceAnalyzeDocuments(t *testing.T) {
	service := NewAnalyzeProposalService(newMockQueueProducer())

	tests := []struct {
		name     string
		cpf      string
		fullName string
		want     bool
	}{
		{name: "empty name", cpf: "12345678902", fullName: "", want: false},
		{name: "minimal valid name", cpf: "12345678901", fullName: "Joe", want: true},
		{name: "valid documents", cpf: "12345678902", fullName: "John Doe", want: true},
		{name: "invalid name too short", cpf: "12345678902", fullName: "Jo", want: false},
		{name: "invalid CPF too short", cpf: "123456789", fullName: "John Doe", want: false},
		{name: "invalid CPF too long", cpf: "123456789012", fullName: "John Doe", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := validPayload()
			payload.CPF = tt.cpf
			payload.FullName = tt.fullName
			got := service.analyzeDocuments(payload)
			if got != tt.want {
				t.Errorf("analyzeDocuments(CPF=%q, FullName=%q) = %v, want %v", tt.cpf, tt.fullName, got, tt.want)
			}
		})
	}
}

func TestAnalyzeProposalServiceAnalyzeCredit(t *testing.T) {
	service := NewAnalyzeProposalService(newMockQueueProducer())

	tests := []struct {
		name   string
		salary float64
		want   bool
	}{
		{name: "above threshold", salary: 5000.0, want: true},
		{name: "at threshold", salary: 3000.0, want: false},
		{name: "below threshold", salary: 2999.99, want: false},
		{name: "just above threshold", salary: 3000.01, want: true},
		{name: "zero", salary: 0.0, want: false},
		{name: "negative", salary: -1.0, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := validPayload()
			payload.Salary = tt.salary
			got := service.analyzeCredit(payload)
			if got != tt.want {
				t.Errorf("analyzeCredit(salary=%.2f) = %v, want %v", tt.salary, got, tt.want)
			}
		})
	}
}

func TestAnalyzeProposalServiceAnalyzeFraud(t *testing.T) {
	service := NewAnalyzeProposalService(newMockQueueProducer())

	tests := []struct {
		name string
		cpf  string
		want bool
	}{
		{name: "even digit", cpf: "12345678902", want: true},
		{name: "odd digit", cpf: "12345678901", want: false},
		{name: "zero", cpf: "12345678900", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := validPayload()
			payload.CPF = tt.cpf
			got := service.analyzeFraud(payload)
			if got != tt.want {
				t.Errorf("analyzeFraud(cpf=%q) = %v, want %v", tt.cpf, got, tt.want)
			}
		})
	}
}

func TestAnalyzeProposalServiceHandle(t *testing.T) {
	tests := []struct {
		name           string
		payload        *entities.ProposalPayload
		wantEvents     int
		wantEventTypes []string
		wantApproved   []bool
	}{
		{
			name:           "documents rejection",
			payload:        &entities.ProposalPayload{CPF: "123", FullName: "John Doe", Salary: 5000.0},
			wantEvents:     1,
			wantEventTypes: []string{events.EventDocumentsRejected},
			wantApproved:   []bool{false},
		},
		{
			name:           "credit rejection",
			payload:        &entities.ProposalPayload{CPF: "12345678902", FullName: "John Doe", Salary: 2000.0},
			wantEvents:     2,
			wantEventTypes: []string{events.EventDocumentsApproved, events.EventCreditRejected},
			wantApproved:   []bool{true, false},
		},
		{
			name:           "fraud rejection",
			payload:        &entities.ProposalPayload{CPF: "12345678901", FullName: "John Doe", Salary: 5000.0},
			wantEvents:     2,
			wantEventTypes: []string{events.EventDocumentsApproved, events.EventFraudRejected},
			wantApproved:   []bool{true, false},
		},
		{
			name:           "all approved",
			payload:        &entities.ProposalPayload{CPF: "12345678902", FullName: "John Doe", Salary: 5000.0},
			wantEvents:     2,
			wantEventTypes: []string{events.EventDocumentsApproved, events.EventRiskAnalysisCompleted},
			wantApproved:   []bool{true, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockQueueProducer()
			service := NewAnalyzeProposalService(mock)
			ctx := context.Background()
			proposalID := uuid.New()

			payloadBytes, _ := json.Marshal(tt.payload)

			event := &entities.IncomingEvent{
				EventType:  events.EventProposalCreated,
				ProposalID: proposalID,
				Payload:    payloadBytes,
			}

			_ = service.Handle(ctx, event)

			assertEventCount(t, mock.published, tt.wantEvents)
			for i := 0; i < tt.wantEvents; i++ {
				assertEvent(t, mock.published[i], tt.wantEventTypes[i], tt.wantApproved[i], proposalID)
			}
		})
	}
}

func TestAnalyzeProposalServiceHandleInvalidPayloadError(t *testing.T) {
	mock := newMockQueueProducer()
	service := NewAnalyzeProposalService(mock)
	ctx := context.Background()
	expectedEventCount := 0

	event := &entities.IncomingEvent{
		EventType:  events.EventProposalCreated,
		ProposalID: uuid.New(),
		Payload:    []byte("invalid json"),
	}

	_ = service.Handle(ctx, event)
	assertEventCount(t, mock.published, expectedEventCount)
}
