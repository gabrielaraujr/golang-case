package ports

import (
	"context"

	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain/events"
)

type QueueProducer interface {
	Publish(ctx context.Context, event *events.RiskEvent) error
}
