package ports

import (
	"context"

	events "github.com/gabrielaraujr/golang-case/account/internal/domain"
)

type QueueProducer interface {
	Publish(ctx context.Context, event *events.ProposalCreatedEvent) error
}
