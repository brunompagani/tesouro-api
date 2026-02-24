.PHONY: help build run test clean docker-build docker-run docker-test docker-clean

# Default target
help:
	@echo "Available targets:"
	@echo "  make build          - Build the Go binary locally"
	@echo "  make run            - Run the updater locally"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run updater in Docker"
	@echo "  make docker-test    - Run tests in Docker"
	@echo "  make all            - Build and test everything"

# Local development
build:
	@echo "Building binary..."
	@go build -o update ./cmd/update

run: build
	@echo "Running updater..."
	@./update

test:
	@echo "Running tests..."
	@go test -v ./cmd/update

clean:
	@echo "Cleaning build artifacts..."
	@rm -f update
	@rm -rf public/*.tmp

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t tesouro-updater:latest .

docker-run: docker-build
	@echo "Running updater in Docker..."
	@docker-compose run --rm updater

docker-test:
	@echo "Running tests in Docker..."
	@docker build -t tesouro-updater:test -f Dockerfile.test .
	@docker run --rm tesouro-updater:test

# Convenience targets
all: test docker-build
	@echo "All done!"

# Development workflow
dev: clean test build run
