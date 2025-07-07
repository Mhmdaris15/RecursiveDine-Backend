## Development Commands

### Setup
```bash
# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Update database credentials in .env file
# Then run migrations
go run cmd/migrate/main.go

# Seed initial data
go run cmd/seed/main.go
```

### Running the application
```bash
# Development mode
go run cmd/api/main.go

# Build and run
go build -o bin/api cmd/api/main.go
./bin/api
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test
go test -v ./internal/services -run TestAuthService_Register
```

### Database Operations
```bash
# Run migrations
go run cmd/migrate/main.go

# Seed sample data
go run cmd/seed/main.go

# Connect to database
psql -h localhost -U postgres -d recursive_dine
```

### Docker
```bash
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
```bash
# Generate Swagger documentation
swag init -g cmd/api/main.go

# Access Swagger UI
# http://localhost:8080/swagger/index.html
```

### Linting
```bash
# Run linter
golangci-lint run

# Auto-fix issues
golangci-lint run --fix
```

### Production Deployment
```bash
# Build production image
docker build -t recursive-dine-api .

# Run production container
docker run -p 8080:8080 --env-file .env recursive-dine-api
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - User logout

### Tables
- `GET /api/v1/tables/{qr_code}` - Get table by QR code

### Menu
- `GET /api/v1/menu` - Get complete menu
- `GET /api/v1/menu/categories` - Get menu categories
- `GET /api/v1/menu/items/search?q={query}` - Search menu items

### Orders
- `POST /api/v1/orders` - Create new order
- `GET /api/v1/orders` - Get user orders
- `GET /api/v1/orders/{id}` - Get order details
- `PATCH /api/v1/orders/{id}/status` - Update order status (staff/admin)

### Payments
- `POST /api/v1/payments/qris` - Initiate QRIS payment
- `POST /api/v1/payments/verify` - Verify payment
- `GET /api/v1/payments/status/{payment_id}` - Get payment status

### Kitchen WebSocket
- `WS /kitchen/updates?token={jwt_token}` - Real-time kitchen updates

### Monitoring
- `GET /metrics` - Prometheus metrics
- `GET /health` - Health check

## Environment Variables

Required environment variables:

- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `JWT_SECRET` - JWT secret key
- `REDIS_HOST` - Redis host
- `REDIS_PORT` - Redis port
- `QRIS_MERCHANT_ID` - QRIS merchant ID
- `QRIS_SECRET_KEY` - QRIS secret key

## Security Features

- JWT authentication with refresh tokens
- Role-based access control
- Rate limiting (100 requests/minute/IP)
- CORS configuration
- Input validation
- SQL injection prevention
- Password hashing with bcrypt
- Secure headers middleware

## Performance Optimizations

- Database connection pooling
- Redis caching
- Optimized database queries with indexes
- Structured logging
- Prometheus metrics collection
- WebSocket for real-time updates

## Testing Strategy

- Unit tests for all services
- Integration tests for API endpoints
- Mock repositories for testing
- Test database setup
- Coverage reporting

## Deployment

The application is containerized and can be deployed using:

1. **Docker Compose** (development)
2. **Kubernetes** (production)
3. **Cloud services** (AWS, GCP, Azure)

## Monitoring

- Prometheus metrics collection
- Grafana dashboards
- Structured logging with logrus
- Health check endpoints
- Error tracking and monitoring
