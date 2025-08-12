@echo off
echo 🧹 Cleaning up existing Docker containers and volumes...

REM Stop all containers
docker-compose -f docker-compose.yml down -v

REM Remove any existing containers
docker container prune -f

REM Remove volumes (this will recreate the database from scratch)
docker volume rm recursivedine-backend_postgres_data 2>nul

echo 🚀 Starting fresh containers...

REM Start the containers
docker-compose -f docker-compose.yml up -d

echo 📊 Checking container status...
docker-compose -f docker-compose.yml ps

echo 📝 Follow logs with: docker-compose -f docker-compose.yml logs -f
echo 🔍 Check database with: docker exec -it recursivedine-backend-postgres-1 psql -U postgres -d recursive_dine
