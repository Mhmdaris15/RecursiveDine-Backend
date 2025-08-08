#!/bin/bash

# Development setup script for RecursiveDine

echo "Starting RecursiveDine Development Environment..."

# Stop any existing containers
echo "Stopping existing containers..."
docker-compose -f docker-compose.dev.yml down 2>/dev/null

# Start development services (database, redis, prometheus, grafana)
echo "Starting development services..."
docker-compose -f docker-compose.dev.yml up -d

# Wait for services to be ready
echo "Waiting for services to start..."
sleep 10

# Check if PostgreSQL is ready
echo "Checking PostgreSQL connection..."
while ! docker exec recursivedine-postgres-dev pg_isready -U postgres -d recursive_dine >/dev/null 2>&1; do
    echo "Waiting for PostgreSQL to be ready..."
    sleep 2
done

echo "PostgreSQL is ready!"

# Run migrations
echo "Running database migrations..."
go run cmd/migrate/migrate.go up

echo ""
echo "======================================"
echo "Development environment is ready!"
echo "======================================"
echo ""
echo "Services running:"
echo "- PostgreSQL: localhost:5432"
echo "- Redis: localhost:6379"
echo "- Prometheus: http://localhost:9090"
echo "- Grafana: http://localhost:3000 (admin/admin)"
echo ""
echo "To start your Go application:"
echo "  go run cmd/api/main.go"
echo ""
echo "To check migration status:"
echo "  go run cmd/migrate/migrate.go status"
echo ""
echo "To stop services:"
echo "  docker-compose -f docker-compose.dev.yml down"
