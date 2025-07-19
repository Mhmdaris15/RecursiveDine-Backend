.PHONY: build run test clean docker-build docker-up docker-down

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=recursive-dine-api
BINARY_PATH=./bin/$(BINARY_NAME)

# Build the application
build:
	$(GOBUILD) -o $(BINARY_PATH) -v ./cmd/api

# Run the application
run:
	$(GOBUILD) -o $(BINARY_PATH) -v ./cmd/api
	$(BINARY_PATH)

# Test the application
test:
	$(GOTEST) -v ./...

# Test with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)
	rm -f coverage.out coverage.html

# Tidy dependencies
deps:
	$(GOMOD) tidy
	$(GOMOD) download

# Run database migrations
migrate:
	$(GOCMD) run cmd/migrate/main.go

# Seed database
seed:
	$(GOCMD) run cmd/seed/main.go

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Development setup
setup: deps migrate seed
	@echo "Development environment setup complete!"

# Quick start for development
dev: setup run

# Production build
build-prod:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo -o $(BINARY_PATH) ./cmd/api
