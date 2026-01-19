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
	if len(args) > 0 {
		log.Printf("[INFO] %s %v", msg, args)
	} else {
		log.Printf("[INFO] %s", msg)
	}
}

func (l *SimpleLogger) Error(ctx context.Context, msg string, args ...any) {
	if len(args) > 0 {
		log.Printf("[ERROR] %s %v", msg, args)
	} else {
		log.Printf("[ERROR] %s", msg)
	}
}

func (l *SimpleLogger) Warn(ctx context.Context, msg string, args ...any) {
	if len(args) > 0 {
		log.Printf("[WARN] %s %v", msg, args)
	} else {
		log.Printf("[WARN] %s", msg)
	}
}
