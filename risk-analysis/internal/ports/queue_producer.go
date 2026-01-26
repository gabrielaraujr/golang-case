package ports

import (
	"context"

	events "github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain"
)

type QueueProducer interface {
	Publish(ctx context.Context, event *events.RiskEvent) error
}
