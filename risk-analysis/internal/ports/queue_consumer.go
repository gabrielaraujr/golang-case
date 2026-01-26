package ports

import (
	"context"

	events "github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain"
)

type EventHandler interface {
	Handle(ctx context.Context, event *events.ProposalCreatedEvent) error
}

type QueueConsumer interface {
	Start(ctx context.Context) error
	Stop() error
}
