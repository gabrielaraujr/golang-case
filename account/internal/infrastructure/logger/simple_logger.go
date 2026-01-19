package logger

import (
	"context"
	"log"
)

type SimpleLogger struct{}

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{}
}

func (l *SimpleLogger) Info(ctx context.Context, msg string, args ...any) {
	log.Printf("[INFO] "+msg, args...)
}

func (l *SimpleLogger) Error(ctx context.Context, msg string, args ...any) {
	log.Printf("[ERROR] "+msg, args...)
}

func (l *SimpleLogger) Warn(ctx context.Context, msg string, args ...any) {
	log.Printf("[WARN] "+msg, args...)
}
