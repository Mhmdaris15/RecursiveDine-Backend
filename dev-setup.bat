@echo off
REM Development setup script for RecursiveDine

echo Starting RecursiveDine Development Environment...

REM Stop any existing containers
echo Stopping existing containers...
docker-compose -f docker-compose.dev.yml down 2>NUL

REM Start development services (database, redis, prometheus, grafana)
echo Starting development services...
docker-compose -f docker-compose.dev.yml up -d

REM Wait for services to be ready
echo Waiting for services to start...
timeout /t 10 /nobreak > NUL

REM Check if PostgreSQL is ready
echo Checking PostgreSQL connection...
:POSTGRES_CHECK
docker exec recursivedine-postgres-dev pg_isready -U postgres -d recursivedine > NUL 2>&1
if %ERRORLEVEL% neq 0 (
    echo Waiting for PostgreSQL to be ready...
    timeout /t 2 /nobreak > NUL
    goto POSTGRES_CHECK
)

echo PostgreSQL is ready!

REM Run migrations
echo Running database migrations...
go run cmd/migrate/migrate.go up

echo.
echo ======================================
echo Development environment is ready!
echo ======================================
echo.
echo Services running:
echo - PostgreSQL: localhost:5432
echo - Redis: localhost:6379  
echo - Prometheus: http://localhost:9090
echo - Grafana: http://localhost:3000 (admin/admin)
echo.
echo To start your Go application:
echo   go run cmd/api/main.go
echo.
echo To check migration status:
echo   go run cmd/migrate/migrate.go status
echo.
echo To stop services:
echo   docker-compose -f docker-compose.dev.yml down
