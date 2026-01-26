package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	events "github.com/gabrielaraujr/golang-case/account/internal/domain"
)

type SQSConfig struct {
	QueueURL string
}

type SQSProducer struct {
	queueURL string
}

func NewSQSProducer(cfg SQSConfig) (*SQSProducer, error) {
	return &SQSProducer{queueURL: cfg.QueueURL}, nil
}

func (p *SQSProducer) Publish(ctx context.Context, event *events.ProposalCreatedEvent) error {
	body, _ := json.Marshal(event)

	form := url.Values{
		"Action":      {"SendMessage"},
		"MessageBody": {string(body)},
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", p.queueURL, bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
