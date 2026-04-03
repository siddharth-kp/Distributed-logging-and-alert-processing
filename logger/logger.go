package logger

import (
	"fmt"
	"log/slog"
)

var logger *slog.Logger

type Backend int

const (
	BackendRedis Backend = iota
	BackendKafka
)

func InitLogger(serviceName string, backend Backend) {
	fmt.Printf("Initializing logger for %s (backend=%d)...\n", serviceName, backend)

	var handler slog.Handler
	switch backend {
	case BackendRedis:
		handler = NewRedisHandler(serviceName)
	case BackendKafka:
		panic("Kafka backend not implemented")
	default:
		panic(fmt.Sprintf("unsupported logger backend: %d", backend))
	}

	logger = slog.New(handler)
}

func GetLogger() *slog.Logger {
	fmt.Println("Getting logger instance...", logger)
	return logger
}
