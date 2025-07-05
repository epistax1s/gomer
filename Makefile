.PHONY: build run test clean docker-build docker-run dev-setup

# Build the application
build:
	go build -o bin/gomer ./cmd/gomer/main.go

# Run the application locally
run: build
	./bin/gomer

# Run in development mode with environment variables
dev:
	. scripts/load-env.sh && go run ./cmd/gomer/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t gomer:latest .

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose
docker-stop:
	docker-compose down

# Setup development environment
dev-setup:
	cp env.example .env
	@echo "Development environment setup complete!"
	@echo "Please edit .env file with your configuration values."

# Show logs
logs:
	docker-compose logs -f

# Database operations
db-migrate:
	go run ./cmd/gomer/main.go migrate

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  dev          - Run in development mode"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker Compose"
	@echo "  dev-setup    - Setup development environment"
	@echo "  logs         - Show Docker logs"
	@echo "  db-migrate   - Run database migrations"
	@echo "  help         - Show this help" 