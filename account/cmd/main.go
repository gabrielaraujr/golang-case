package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpRouter "github.com/gabrielaraujr/golang-case/account/internal/adapters/http"
	"github.com/gabrielaraujr/golang-case/account/internal/adapters/http/handler"
	"github.com/gabrielaraujr/golang-case/account/internal/application/services"
	"github.com/gabrielaraujr/golang-case/account/internal/infrastructure/logger"
	"github.com/gabrielaraujr/golang-case/account/internal/infrastructure/postgres"
	"github.com/gabrielaraujr/golang-case/account/internal/infrastructure/queue"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	queueProducer, err := queue.NewSQSProducer(queue.SQSConfig{
		QueueURL: os.Getenv("SQS_QUEUE_URL"),
	})
	if err != nil {
		log.Fatalf("Failed to initialize SQS producer: %v", err)
	}

	proposalRepo := postgres.NewProposalRepository(dbPool)
	createUC := services.NewCreateProposalUseCase(proposalRepo, queueProducer, logger.NewSimpleLogger())
	getUC := services.NewGetProposalUseCase(proposalRepo)

	router := httpRouter.NewRouter(handler.NewProposalHandler(createUC, getUC))

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", os.Getenv("PORT"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Printf("server exited")
}
