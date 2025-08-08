@echo off
REM Development runner script - starts the Go application locally

REM Check if .env.dev exists
if not exist ".env.dev" (
    echo Error: .env.dev file not found!
    echo Please copy .env.dev.example to .env.dev and configure it.
    exit /b 1
)

REM Check if services are running
echo Checking development services...
docker ps | findstr "recursivedine-postgres-dev" > NUL
if %ERRORLEVEL% neq 0 (
    echo Error: Development services are not running!
    echo Please run: dev-setup.bat
    exit /b 1
)

echo Starting RecursiveDine API server...
echo.
echo Server will start on: http://localhost:8002
echo Press Ctrl+C to stop the server
echo.

REM Set environment to development
set APP_ENV=development

REM Run the Go application with the development environment
go run cmd/api/main.go
