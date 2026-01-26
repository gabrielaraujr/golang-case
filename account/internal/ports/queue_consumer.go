package ports

import (
	"context"

	events "github.com/gabrielaraujr/golang-case/account/internal/domain"
)

type EventHandler interface {
	Handle(ctx context.Context, event *events.ProposalStatusChangedEvent) error
}
