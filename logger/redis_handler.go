package logger

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHandler struct {
	rdb         *redis.Client
	serviceName string
}

func NewRedisHandler(serviceName string) slog.Handler {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	return &RedisHandler{
		rdb:         rdb,
		serviceName: serviceName,
	}
}

func (h *RedisHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *RedisHandler) Handle(ctx context.Context, record slog.Record) error {
	logEntry := map[string]interface{}{
		"service_name": h.serviceName,
		"level":        record.Level.String(),
		"message":      record.Message,
		"timestamp":    time.Now(),
	}

	jsonLog, _ := json.Marshal(logEntry)

	// async push (important)
	go h.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "logs_stream",
		Values: map[string]interface{}{
			"data": jsonLog,
		},
	})

	return nil
}

func (h *RedisHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *RedisHandler) WithGroup(name string) slog.Handler {
	return h
}
