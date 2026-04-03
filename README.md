# DLPAS - Distributed Log Processing and Alerting System

A Go-based distributed system for collecting, processing, and alerting on logs from multiple services using Redis streams.

## Architecture

- **Service 1 & 2**: Log producers that generate logs continuously and send them to Redis streams.
- **Collector**: Consumes logs from Redis, enriches them with timestamps, and processes them (e.g., for alerting).
- **Logger Package**: Configurable logger with backends (currently Redis, extensible to Kafka).
- **Redis**: Acts as the message queue for logs using streams (`logs_stream` and `enriched_logs_stream`).

## Prerequisites

- Go 1.19+
- Redis (running locally or in Docker)
- Docker (optional, for containerized Redis)

## Installation

1. Clone the repository:
   ```bash
   git clone <your-repo-url>
   cd DLPAS
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Start Redis:
   - Using Docker:
     ```bash
     docker run -d -p 6379:6379 --name redis-server redis
     ```
   - Or install Redis locally.

## Usage

### Running the Full System

1. Set environment variable:
   ```bash
   export REDIS_ADDR=localhost:6379
   ```

2. Start all services:
   ```bash
   go run main.go
   ```
   This launches Service 1 (port 8081), Service 2 (port 8082), and the Collector.

### Running Individual Components

- **Service 1**:
  ```bash
  go run service1/main.go
  ```

- **Service 2**:
  ```bash
  go run service2/main.go
  ```

- **Collector**:
  ```bash
  go run Collector/main.go
  ```

### Testing Logs

1. Hit the endpoints to generate logs:
   - Service 1: `curl http://localhost:8081/`
   - Service 2: `curl http://localhost:8083/`

2. Check logs in Redis:
   - View stream length: `docker exec -it redis-server redis-cli XLEN logs_stream`
   - Read logs: `docker exec -it redis-server redis-cli XRANGE logs_stream - +`
   - View enriched logs: `docker exec -it redis-server redis-cli XRANGE enriched_logs_stream - +`

3. Clear logs: `docker exec -it redis-server redis-cli DEL logs_stream`

## Logger Configuration

The logger supports backends:
- `0`: Redis (default)
- `1`: Kafka (not implemented yet)

Example:
```go
logger.InitLogger("service1", logger.BackendRedis)
```

## Development

- Add tests: `go test ./...`
- Build: `go build`
- Format code: `go fmt ./...`

## Next Steps

- Implement Kafka backend in logger.
- Add alerting logic in `handleEnrichedLog()`.
- Containerize with Docker Compose.
- Add monitoring and metrics.

## License

MIT