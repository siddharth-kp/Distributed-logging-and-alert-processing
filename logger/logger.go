package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

var logger *slog.Logger

func InitLogger(serviceName string) {
	fmt.Printf("Initializing logger for %s...\n", serviceName)
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	handler := &RedisHandler{
		rdb:         rdb,
		serviceName: serviceName,
	}

	logger = slog.New(handler)
}

func GetLogger() *slog.Logger {
	fmt.Println("Getting logger instance...", logger)
	return logger
}