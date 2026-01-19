package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain/events"
)

type SQSConfig struct {
	QueueURL string
}

type SQSProducer struct {
	queueURL string
	client   *http.Client
}

func NewSQSProducer(cfg SQSConfig) (*SQSProducer, error) {
	if cfg.QueueURL == "" {
		return nil, fmt.Errorf("SQS_PROPOSALS_QUEUE_URL is required")
	}

	return &SQSProducer{
		queueURL: cfg.QueueURL,
		client:   &http.Client{},
	}, nil
}

func (p *SQSProducer) Publish(ctx context.Context, event *events.RiskEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	form := url.Values{
		"Action":      {"SendMessage"},
		"MessageBody": {string(body)},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.queueURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sqs error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
