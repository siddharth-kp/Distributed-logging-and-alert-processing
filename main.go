// Distributed Log Processing and Alerting System
// Architecture:
//   - Service 1 & 2: Log Producers (generate logs continuously)
//   - Service 3: Log Processor (consumes, processes, alerts)
//   - Queue: In-memory log queue (will upgrade to Redis/Kafka later)

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Starting DLPAS - Distributed Log Processing and Alerting System\n")

	// Start Service 1 (Log Producer)
	service1Cmd := exec.Command("go", "run", "./service1/main.go")
	service1Cmd.Stdout = os.Stdout
	service1Cmd.Stderr = os.Stderr
	if err := service1Cmd.Start(); err != nil {
		log.Fatalf("Failed to start service1: %v", err)
	}
	fmt.Println("✓ Service 1 (Producer) started on :8081")

	// Start Service 2 (Log Producer)
	service2Cmd := exec.Command("go", "run", "./service2/main.go")
	service2Cmd.Stdout = os.Stdout
	service2Cmd.Stderr = os.Stderr
	if err := service2Cmd.Start(); err != nil {
		log.Fatalf("Failed to start service2: %v", err)
	}
	fmt.Println("✓ Service 2 (Producer) started on :8082")

	// Start Service 3 (Log Processor/Consumer)
	service3Cmd := exec.Command("go", "run", "./service3/main.go")
	service3Cmd.Stdout = os.Stdout
	service3Cmd.Stderr = os.Stderr
	if err := service3Cmd.Start(); err != nil {
		log.Fatalf("Failed to start service3: %v", err)
	}
	fmt.Println("✓ Service 3 (Processor) started")

	fmt.Println("\nAll services running. Press Ctrl+C to stop.\n")

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nShutting down services...")
	service1Cmd.Process.Kill()
	service2Cmd.Process.Kill()
	// // service3Cmd.Process.Kill()
	fmt.Println("✓ All services stopped")
}
