package entities

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ProposalPayload struct {
	FullName string  `json:"full_name"`
	CPF      string  `json:"cpf"`
	Salary   float64 `json:"salary"`
}

type IncomingEvent struct {
	EventType  string          `json:"event_type"`
	ProposalID uuid.UUID       `json:"proposal_id"`
	Payload    json.RawMessage `json:"payload"`
}

func (e *IncomingEvent) ParsePayload() (*ProposalPayload, error) {
	var payload ProposalPayload
	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
