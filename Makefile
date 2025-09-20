.PHONY: build run test clean dev fmt lint install-tools

# Variables
BINARY_NAME=hex-arch-server
MAIN_PATH=cmd/server/main.go
PORT?=8080
WEATHER_API_KEY?=demo-key

# Build the application
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	PORT=$(PORT) WEATHER_API_KEY=$(WEATHER_API_KEY) go run $(MAIN_PATH)

# Run with hot reload (requires air)
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	PORT=$(PORT) WEATHER_API_KEY=$(WEATHER_API_KEY) air

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	go fmt ./...
	gofmt -w .

# Run linter (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Clean build artifacts
clean:
	go clean
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install development tools
install-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker build
docker-build:
	docker build -t $(BINARY_NAME) .

# Docker run
docker-run:
	docker run -p $(PORT):$(PORT) -e PORT=$(PORT) -e WEATHER_API_KEY=$(WEATHER_API_KEY) $(BINARY_NAME)

# Help
help:
	@echo "Available targets:"
	@echo "  make build         - Build the binary"
	@echo "  make run           - Run the application"
	@echo "  make dev           - Run with hot reload (requires air)"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
	@echo "  make tidy          - Tidy go modules"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make install-tools - Install development tools"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
	@echo ""
	@echo "Environment variables:"
	@echo "  PORT              - Server port (default: 8080)"
	@echo "  WEATHER_API_KEY   - OpenWeather API key (default: demo-key)"