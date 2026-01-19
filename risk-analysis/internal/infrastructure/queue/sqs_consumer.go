package queue

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/domain/entities"
	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/ports"
)

type SQSConsumerConfig struct {
	QueueURL        string
	PollingInterval time.Duration
	MaxMessages     int
}

type SQSConsumer struct {
	queueURL        string
	pollingInterval time.Duration
	maxMessages     int
	client          *http.Client
	handler         ports.EventHandler
	stopCh          chan struct{}
	wg              sync.WaitGroup
	running         bool
	mu              sync.Mutex
}

func NewSQSConsumer(cfg SQSConsumerConfig, handler ports.EventHandler) (*SQSConsumer, error) {
	if cfg.QueueURL == "" {
		return nil, fmt.Errorf("SQS_PROPOSALS_QUEUE_URL is required")
	}
	if handler == nil {
		return nil, fmt.Errorf("event handler is required")
	}

	pollingInterval := cfg.PollingInterval
	if pollingInterval == 0 {
		pollingInterval = 5 * time.Second
	}

	maxMessages := cfg.MaxMessages
	if maxMessages == 0 {
		maxMessages = 10
	}

	return &SQSConsumer{
		queueURL:        cfg.QueueURL,
		pollingInterval: pollingInterval,
		maxMessages:     maxMessages,
		client:          &http.Client{Timeout: 30 * time.Second},
		handler:         handler,
		stopCh:          make(chan struct{}),
	}, nil
}

func (c *SQSConsumer) Start(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return fmt.Errorf("consumer already running")
	}
	c.running = true
	c.mu.Unlock()

	log.Printf("[SQSConsumer] Starting consumer for queue: %s", c.queueURL)

	c.wg.Add(1)
	go c.pollMessages(ctx)

	return nil
}

func (c *SQSConsumer) Stop() error {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return nil
	}
	c.running = false
	c.mu.Unlock()

	close(c.stopCh)
	c.wg.Wait()

	log.Printf("[SQSConsumer] Consumer stopped")
	return nil
}

func (c *SQSConsumer) pollMessages(ctx context.Context) {
	defer c.wg.Done()

	ticker := time.NewTicker(c.pollingInterval)
	defer ticker.Stop()

	// Poll immediately on start
	c.receiveAndProcess(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[SQSConsumer] Context cancelled, stopping")
			return
		case <-c.stopCh:
			log.Printf("[SQSConsumer] Stop signal received")
			return
		case <-ticker.C:
			c.receiveAndProcess(ctx)
		}
	}
}

func (c *SQSConsumer) receiveAndProcess(ctx context.Context) {
	messages, err := c.receiveMessages(ctx)
	if err != nil {
		log.Printf("[SQSConsumer] Error receiving messages: %v", err)
		return
	}

	for _, msg := range messages {
		if err := c.processMessage(ctx, msg); err != nil {
			log.Printf("[SQSConsumer] Error processing message: %v", err)
			continue
		}

		if err := c.deleteMessage(ctx, msg.ReceiptHandle); err != nil {
			log.Printf("[SQSConsumer] Error deleting message: %v", err)
		}
	}
}

// SQS XML response structures
type receiveMessageResponse struct {
	XMLName xml.Name `xml:"ReceiveMessageResponse"`
	Result  struct {
		Messages []sqsMessage `xml:"Message"`
	} `xml:"ReceiveMessageResult"`
}

type sqsMessage struct {
	MessageId     string `xml:"MessageId"`
	ReceiptHandle string `xml:"ReceiptHandle"`
	Body          string `xml:"Body"`
}

func (c *SQSConsumer) receiveMessages(ctx context.Context) ([]sqsMessage, error) {
	form := url.Values{
		"Action":              {"ReceiveMessage"},
		"MaxNumberOfMessages": {fmt.Sprintf("%d", c.maxMessages)},
		"WaitTimeSeconds":     {"5"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.queueURL+"?"+form.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("receive messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("sqs error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var response receiveMessageResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return response.Result.Messages, nil
}

func (c *SQSConsumer) processMessage(ctx context.Context, msg sqsMessage) error {
	log.Printf("[SQSConsumer] Processing message: %s", msg.MessageId)

	var event entities.IncomingEvent
	if err := json.Unmarshal([]byte(msg.Body), &event); err != nil {
		return fmt.Errorf("unmarshal event: %w", err)
	}

	return c.handler.Handle(ctx, &event)
}

func (c *SQSConsumer) deleteMessage(ctx context.Context, receiptHandle string) error {
	form := url.Values{
		"Action":        {"DeleteMessage"},
		"ReceiptHandle": {receiptHandle},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.queueURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.URL.RawQuery = form.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("delete message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sqs error (status %d): %s", resp.StatusCode, string(body))
	}

	log.Printf("[SQSConsumer] Message deleted successfully")
	return nil
}
