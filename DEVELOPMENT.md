## Development Commands

### Setup
```powershell
# Install dependencies
go mod download

# Copy environment file
Copy-Item .env.example .env

# Update database credentials in .env file
# Then run migrations
go run cmd/migrate/main.go

# Seed initial data
go run cmd/seed/main.go
```

### Running the application
```powershell
# Development mode
go run cmd/api/main.go

# Build and run
go build -o bin/api.exe cmd/api/main.go
./bin/api.exe
```

### Testing
```powershell
# Run all tests
go test ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test
go test -v ./internal/services -run TestAuthService_Register
```

### Database Operations
```powershell
# Run migrations
go run cmd/migrate/main.go

# Seed sample data
go run cmd/seed/main.go

# Connect to database (if psql is installed)
psql -h localhost -U postgres -d recursive_dine
```

### Docker
```powershell
# Build and run with Docker Compose
docker-compose up --build

# Run only database services
docker-compose up postgres redis

# Stop services
docker-compose down

# Remove volumes
docker-compose down -v
```

### API Documentation
```powershell
# Generate Swagger documentation (if swag is installed)
swag init -g cmd/api/main.go

# Access Swagger UI
# http://localhost:8080/swagger/index.html
```

### Linting
```powershell
# Run linter (if golangci-lint is installed)
golangci-lint run

# Auto-fix issues
golangci-lint run --fix
```

### Production Deployment
```powershell
# Build production image
docker build -t recursive-dine-api .

# Run production container
docker run -p 8080:8080 --env-file .env recursive-dine-api
```

## Installation Requirements

### Prerequisites
1. **Go 1.24+** - [Download](https://golang.org/dl/)
2. **PostgreSQL 15+** - [Download](https://www.postgresql.org/download/)
3. **Redis** (optional) - [Download](https://redis.io/download/)
4. **Docker** (optional) - [Download](https://docs.docker.com/get-docker/)

### Optional Tools
- **golangci-lint** - For code linting
- **swag** - For API documentation generation
- **psql** - PostgreSQL command-line tool

## Quick Start

1. **Clone and setup**:
```powershell
git clone <repository-url>
cd RecursiveDine
go mod download
Copy-Item .env.example .env
```

2. **Configure database**:
   - Edit `.env` file with your database credentials
   - Ensure PostgreSQL is running

3. **Initialize database**:
```powershell
go run cmd/migrate/main.go
go run cmd/seed/main.go
```

4. **Run the application**:
```powershell
go run cmd/api/main.go
```

5. **Test the API**:
   - Health check: `http://localhost:8080/health`
   - API documentation: `http://localhost:8080/swagger/index.html`

## Default Credentials

After seeding the database, you can use these default accounts:

- **Admin**: `admin` / `admin123`
- **Staff**: `staff1` / `staff123`
- **Customer**: `customer1` / `customer123`

## API Testing

You can use tools like:
- **Postman** - Import the API collection
- **curl** - Command-line HTTP client
- **Insomnia** - REST client
- **Thunder Client** - VS Code extension

Example API calls:
```powershell
# Register new user
curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Get menu
curl http://localhost:8080/api/v1/menu

# Get table info
curl http://localhost:8080/api/v1/tables/QR001
```
