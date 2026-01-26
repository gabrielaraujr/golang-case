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

	events "github.com/gabrielaraujr/golang-case/account/internal/domain"
	"github.com/gabrielaraujr/golang-case/account/internal/ports"
)

type SQSConsumerConfig struct {
	QueueURL    string
	MaxMessages int
}

type SQSConsumer struct {
	queueURL string
	handler  ports.EventHandler
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

func NewSQSConsumer(cfg SQSConsumerConfig, handler ports.EventHandler) (*SQSConsumer, error) {
	return &SQSConsumer{
		queueURL: cfg.QueueURL,
		handler:  handler,
		stopCh:   make(chan struct{}),
	}, nil
}

func (c *SQSConsumer) Start(ctx context.Context) error {
	c.wg.Add(1)
	go c.poll(ctx)
	return nil
}

func (c *SQSConsumer) Stop() error {
	close(c.stopCh)
	c.wg.Wait()
	return nil
}

func (c *SQSConsumer) poll(ctx context.Context) {
	defer c.wg.Done()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.processMessages(ctx)
		}
	}
}

func (c *SQSConsumer) processMessages(ctx context.Context) {
	messages := c.receive(ctx)
	for _, msg := range messages {
		var event events.ProposalStatusChangedEvent
		if err := json.Unmarshal([]byte(msg.Body), &event); err != nil {
			log.Printf("[SQSConsumer] Error parsing message: %v", err)
			continue
		}

		log.Printf("[SQSConsumer] Processing message: %s", msg.MessageId)
		if err := c.handler.Handle(ctx, &event); err != nil {
			log.Printf("[SQSConsumer] Error processing message: %v", err)
			continue
		}

		c.delete(ctx, msg.ReceiptHandle)
		log.Printf("[SQSConsumer] Message deleted successfully")
	}
}

type message struct {
	MessageId     string `xml:"MessageId"`
	ReceiptHandle string `xml:"ReceiptHandle"`
	Body          string `xml:"Body"`
}

func (c *SQSConsumer) receive(ctx context.Context) []message {
	url := fmt.Sprintf("%s?Action=ReceiveMessage&MaxNumberOfMessages=10&WaitTimeSeconds=5", c.queueURL)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		XMLName xml.Name `xml:"ReceiveMessageResponse"`
		Result  struct {
			Messages []message `xml:"Message"`
		} `xml:"ReceiveMessageResult"`
	}
	xml.Unmarshal(body, &result)
	return result.Result.Messages
}

func (c *SQSConsumer) delete(ctx context.Context, receiptHandle string) {
	form := url.Values{"Action": {"DeleteMessage"}, "ReceiptHandle": {receiptHandle}}
	req, _ := http.NewRequestWithContext(ctx, "POST", c.queueURL+"?"+form.Encode(), nil)
	http.DefaultClient.Do(req)
}
