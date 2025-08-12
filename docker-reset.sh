#!/bin/bash

echo "🧹 Cleaning up existing Docker containers and volumes..."

# Stop all containers
docker-compose -f docker-compose.yml down -v

# Remove any existing containers
docker container prune -f

# Remove volumes (this will recreate the database from scratch)
docker volume rm recursivedine-backend_postgres_data 2>/dev/null || true

echo "🚀 Starting fresh containers..."

# Start the containers
docker-compose -f docker-compose.yml up -d

echo "📊 Checking container status..."
docker-compose -f docker-compose.yml ps

echo "📝 Follow logs with: docker-compose -f docker-compose.yml logs -f"
echo "🔍 Check database with: docker exec -it recursivedine-backend-postgres-1 psql -U postgres -d recursive_dine"
