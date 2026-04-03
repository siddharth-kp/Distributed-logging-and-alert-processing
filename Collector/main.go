package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx := context.Background()

	fmt.Printf("Collector started: redis=%s\n", redisAddr)

	// Read from stream in real time (blocking) using XREAD
	lastID := "$"
	for {
		streamRes, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{"logs_stream", lastID},
			Count:   50,
			Block:   500 * time.Millisecond,
		}).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			fmt.Fprintf(os.Stderr, "error reading stream: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, stream := range streamRes {
			for _, msg := range stream.Messages {
				entryRaw, ok := msg.Values["data"].(string)
				if !ok {
					fmt.Fprintf(os.Stderr, "invalid entry type for message %s\n", msg.ID)
					continue
				}

				var entry map[string]interface{}
				if err := json.Unmarshal([]byte(entryRaw), &entry); err != nil {
					fmt.Fprintf(os.Stderr, "invalid JSON in message %s: %v\n", msg.ID, err)
					continue
				}

				entry["collected_at"] = time.Now().UTC().Format(time.RFC3339Nano)

				enrichedLog, err := json.Marshal(entry)
				if err != nil {
					fmt.Fprintf(os.Stderr, "json marshal failed %s: %v\n", msg.ID, err)
					continue
				}

				_, err = rdb.XAdd(ctx, &redis.XAddArgs{
					Stream: "enriched_logs_stream",
					Values: map[string]interface{}{"data": enrichedLog},
				}).Result()
				if err != nil {
					fmt.Fprintf(os.Stderr, "xadd failed %s: %v\n", msg.ID, err)
					continue
				}

				lastID = msg.ID
				fmt.Printf("enriched message %s\n", msg.ID)

				// callback for further processing
				handleEnrichedLog(entry)
			}
		}
	}
}

func handleEnrichedLog(entry map[string]interface{}) {
	// TODO: process enriched log (alert, write to DB, sink to Kafka, etc.)
	fmt.Println("one line of log received", entry)
}

