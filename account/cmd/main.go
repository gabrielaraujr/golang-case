package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpRouter "github.com/gabrielaraujr/golang-case/account/internal/adapters/http"
	"github.com/gabrielaraujr/golang-case/account/internal/adapters/http/handler"
	"github.com/gabrielaraujr/golang-case/account/internal/application/services"
	"github.com/gabrielaraujr/golang-case/account/internal/infrastructure/logger"
	"github.com/gabrielaraujr/golang-case/account/internal/infrastructure/postgres"
	"github.com/gabrielaraujr/golang-case/account/internal/infrastructure/queue"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Println("[Account] Starting...")

	// Database
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Dependencies
	producer, _ := queue.NewSQSProducer(queue.SQSConfig{
		QueueURL: os.Getenv("SQS_PROPOSALS_QUEUE_URL"),
	})
	repo := postgres.NewProposalRepository(dbPool)
	logger := logger.NewSimpleLogger()

	// Use Cases
	createUC := services.NewCreateProposalUseCase(repo, producer, logger)
	getUC := services.NewGetProposalUseCase(repo)

	// Consumer
	eventHandler := services.NewProposalStatusChangedEventHandler(repo, logger)
	consumer, _ := queue.NewSQSConsumer(queue.SQSConsumerConfig{
		QueueURL:    os.Getenv("SQS_RISK_QUEUE_URL"),
		MaxMessages: 10,
	}, eventHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = consumer.Start(ctx)
	log.Println("[Account] Consumer started")

	// HTTP Server
	port := os.Getenv("PORT")
	router := httpRouter.NewRouter(handler.NewProposalHandler(createUC, getUC))
	go func() {
		log.Printf("[Account] Server listening on :%s", port)
		if err := http.ListenAndServe(":"+port, router); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("[Account] Shutting down...")
	_ = consumer.Stop()
}
