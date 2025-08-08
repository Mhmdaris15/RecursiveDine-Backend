#!/bin/bash

# Development runner script - starts the Go application locally

# Check if .env.dev exists
if [ ! -f ".env.dev" ]; then
    echo "Error: .env.dev file not found!"
    echo "Please copy .env.dev.example to .env.dev and configure it."
    exit 1
fi

# Check if services are running
if ! docker ps | grep -q "recursivedine-postgres-dev"; then
    echo "Error: Development services are not running!"
    echo "Please run: ./dev-setup.sh"
    exit 1
fi

echo "Starting RecursiveDine API server..."
echo ""
echo "Server will start on: http://localhost:8002"
echo "Press Ctrl+C to stop the server"
echo ""

# Set environment to development
export APP_ENV=development

# Run the Go application with the development environment
go run cmd/api/main.go
