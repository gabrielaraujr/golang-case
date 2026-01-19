package domain

import "errors"

var (
	ErrInvalidProposalID = errors.New("invalid proposal id")
	ErrInvalidEventType  = errors.New("invalid event type")
	ErrEmptyPayload      = errors.New("empty payload")
)
