package ports

import "context"

type QueueMessage struct {
	ID         string
	Body       string
	Attributes map[string]string
}

type QueueConsumer interface {
	Start(ctx context.Context, handler func(context.Context, *QueueMessage) error) error
	Close() error
}
