# RecursiveDine Development Environment Setup

This guide explains how to set up a hybrid development environment where database and tools run in Docker while the Go application runs locally for faster development.

## 🏗️ Architecture

**Docker Services (Infrastructure):**
- PostgreSQL Database (port 5432)
- Redis Cache (port 6379)  
- Prometheus Metrics (port 9090)
- Grafana Monitoring (port 3000)

**Local Development:**
- Go API Server (port 8002)
- Hot reloading and debugging
- No Docker rebuild needed

## 🚀 Quick Start

### 1. Start Development Environment

**Windows:**
```bash
dev-setup.bat
```

**Linux/macOS:**
```bash
chmod +x dev-setup.sh
./dev-setup.sh
```

This will:
- Start all Docker services
- Wait for PostgreSQL to be ready
- Run database migrations automatically

### 2. Start Your Go Application

**Windows:**
```bash
dev-run.bat
```

**Linux/macOS:**
```bash
chmod +x dev-run.sh  
./dev-run.sh
```

Your API will be available at: `http://localhost:8002`

## 📊 Monitoring & Tools

Once development environment is running:

- **API Server**: http://localhost:8002
- **PostgreSQL**: localhost:5432 (user: postgres, password: postgres, db: recursive_dine)
- **Redis**: localhost:6379
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

## 🗄️ Database Migration Commands

### Basic Commands

```bash
# Check migration status
go run cmd/migrate/main.go status

# Run pending migrations
go run cmd/migrate/main.go up

# Rollback migrations (1 step)
go run cmd/migrate/main.go down 1

# Create new migration
go run cmd/migrate/main.go create add_new_table

# Reset database (DANGER - deletes all data)
go run cmd/migrate/main.go reset
```

### Migration Workflow

1. **Create a new migration:**
   ```bash
   go run cmd/migrate/main.go create add_payment_methods
   ```
   This creates: `migrations/004_add_payment_methods.sql`

2. **Edit the migration file:**
   ```sql
   -- Migration: add_payment_methods
   -- Created: 2025-08-08 15:30:00

   CREATE TABLE payment_methods (
       id SERIAL PRIMARY KEY,
       name VARCHAR(100) NOT NULL,
       is_active BOOLEAN DEFAULT true,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

   CREATE INDEX idx_payment_methods_name ON payment_methods(name);
   ```

3. **Run the migration:**
   ```bash
   go run cmd/migrate/main.go up
   ```

4. **Check status:**
   ```bash
   go run cmd/migrate/main.go status
   ```

### Example Migration Status Output

```
Migration Status:
================
✓ Applied  001_initial_schema
✓ Applied  002_add_user_fields  
✓ Applied  003_add_vat_and_cashier_fields
✗ Pending  004_add_payment_methods

Total migrations: 4
Applied: 3
Pending: 1
```

## 🔧 Configuration

### Environment Files

- **`.env.dev`** - Development environment (used when APP_ENV=development)
- **`.env`** - Production environment

### Development Environment Variables

```bash
# Database Configuration (Docker)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=recursive_dine
DB_SSL_MODE=disable

# Redis Configuration (Docker)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Application Configuration
APP_ENV=development
APP_PORT=8002
LOG_LEVEL=debug
```

## 🛠️ Development Workflow

### Daily Development

1. **Start development environment once:**
   ```bash
   dev-setup.bat  # Windows
   # OR
   ./dev-setup.sh  # Linux/macOS
   ```

2. **Start your Go app for development:**
   ```bash
   dev-run.bat  # Windows  
   # OR
   ./dev-run.sh  # Linux/macOS
   ```

3. **Make code changes and restart as needed**
   - Just stop (Ctrl+C) and run `dev-run.bat` again
   - No Docker rebuild required!

### Database Changes

1. **Create migration:**
   ```bash
   go run cmd/migrate/main.go create your_change_name
   ```

2. **Edit the SQL file in `migrations/` folder**

3. **Apply migration:**
   ```bash
   go run cmd/migrate/main.go up
   ```

4. **Restart your Go app to see changes**

### Testing

1. **Run tests with development database:**
   ```bash
   go test ./...
   ```

2. **Use Postman collections:**
   - Import `RecursiveDine_E2E_Testing.postman_collection.json`
   - Import `RecursiveDine_E2E_Environment.postman_environment.json`
   - Set base_url to `http://localhost:8002`

## 🐳 Docker Commands

### Manage Development Services

```bash
# Start all services
docker-compose -f docker-compose.dev.yml up -d

# Stop all services  
docker-compose -f docker-compose.dev.yml down

# View logs
docker-compose -f docker-compose.dev.yml logs postgres
docker-compose -f docker-compose.dev.yml logs redis

# Check service status
docker-compose -f docker-compose.dev.yml ps
```

### Database Access

```bash
# Connect to PostgreSQL
docker exec -it recursivedine-postgres-dev psql -U postgres -d recursive_dine

# Connect to Redis
docker exec -it recursivedine-redis-dev redis-cli
```

## 🚨 Troubleshooting

### Common Issues

1. **"Failed to connect to database"**
   ```bash
   # Check if PostgreSQL is running
   docker ps | grep postgres
   
   # Restart if needed
   docker-compose -f docker-compose.dev.yml restart postgres
   ```

2. **"Port already in use"**
   ```bash
   # Check what's using the port
   netstat -ano | findstr :5432  # Windows
   lsof -i :5432                 # Linux/macOS
   
   # Stop conflicting services
   docker-compose -f docker-compose.dev.yml down
   ```

3. **"Migration failed"**
   ```bash
   # Check migration status
   go run cmd/migrate/main.go status
   
   # Rollback if needed
   go run cmd/migrate/main.go down 1
   
   # Fix migration file and retry
   go run cmd/migrate/main.go up
   ```

4. **".env.dev not found"**
   ```bash
   # Copy the example file
   cp .env.dev .env.dev
   
   # Edit configuration as needed
   ```

### Reset Everything

If you need to start fresh:

```bash
# Stop all services
docker-compose -f docker-compose.dev.yml down

# Remove volumes (deletes all data)
docker volume rm recursivedine_postgres_dev_data
docker volume rm recursivedine_redis_dev_data

# Start fresh
dev-setup.bat  # or ./dev-setup.sh
```

## 📝 Benefits of This Setup

✅ **Faster Development**
- No Docker rebuilds when changing Go code
- Instant restarts
- Better debugging experience

✅ **Consistent Environment**  
- Database and tools always consistent
- Easy to reset and recreate
- Isolated from system installs

✅ **Easy Migration Management**
- Simple CLI tool for database changes
- Version control for schema changes
- Easy rollback and status checking

✅ **Professional Monitoring**
- Prometheus metrics out of the box
- Grafana dashboards ready to use
- Redis caching available

This setup gives you the best of both worlds: the convenience of local development with the consistency of containerized infrastructure!
