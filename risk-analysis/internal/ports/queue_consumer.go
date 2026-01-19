package ports

import (
	"context"

	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain/entities"
)

type EventHandler interface {
	Handle(ctx context.Context, event *entities.IncomingEvent) error
}

type QueueConsumer interface {
	Start(ctx context.Context) error
	Stop() error
}
