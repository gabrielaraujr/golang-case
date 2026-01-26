package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/application/services"
	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/infrastructure/logger"
	"github.com/gabrielaraujr/golang-case/risk-analysis/internal/infrastructure/queue"
)

func main() {
	log.Println("[RiskAnalysis] Starting...")

	// Logger
	appLogger := logger.NewSimpleLogger()

	// Producer
	producer, _ := queue.NewSQSProducer(queue.SQSConfig{
		QueueURL: os.Getenv("SQS_RISK_QUEUE_URL"),
	})

	// Service
	analyzeService := services.NewAnalyzeProposalService(producer, appLogger)

	// Consumer
	consumer, _ := queue.NewSQSConsumer(queue.SQSConsumerConfig{
		QueueURL:    os.Getenv("SQS_PROPOSALS_QUEUE_URL"),
		MaxMessages: 10,
	}, analyzeService, appLogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = consumer.Start(ctx)
	log.Println("[RiskAnalysis] Consumer started")

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("[RiskAnalysis] Shutting down...")
	_ = consumer.Stop()
}
