## RecursiveDine API

This is the backend API for RecursiveDine, a restaurant management system.

### Features

- User authentication (JWT)
- Table management
- Menu management
- Order processing
- QRIS payment integration
- Real-time kitchen updates (WebSocket)
- Prometheus monitoring

### Getting Started

#### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make

#### Setup and Running

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd RecursiveDine
   ```

2. **Copy the environment file:**
   ```bash
   cp .env.example .env
   ```
   *Update `.env` with your database credentials.*

3. **Run the development environment:**
   ```bash
   make dev
   ```
   This command will:
   - Install dependencies
   - Run database migrations
   - Seed the database
   - Start the API server

   The API will be available at `http://localhost:8002`.

### Development

The `Makefile` contains all the necessary commands for development.

| Command | Description |
|---|---|
| `make build` | Build the application binary. |
| `make run` | Build and run the application. |
| `make test` | Run all tests. |
| `make test-coverage` | Run tests and generate an HTML coverage report. |
| `make clean` | Remove build files and coverage reports. |
| `make deps` | Tidy and download Go modules. |
| `make migrate` | Run database migrations. |
| `make seed` | Seed the database with initial data. |
| `make setup` | Set up the development environment (deps, migrate, seed). |
| `make dev` | Run the complete development environment. |
| `make docker-build` | Build the Docker image. |
| `make docker-up` | Start the services using Docker Compose. |
| `make docker-down` | Stop the services. |
| `make docker-logs` | View the logs of the running services. |

### API Documentation

The API is documented using Swagger. Once the application is running, you can access the Swagger UI at:

`http://localhost:8002/swagger/index.html`

To generate the Swagger documentation, run:
```bash
swag init -g cmd/api/main.go
```

### Default Credentials

After seeding the database, you can use these default accounts:

- **Admin**: `admin` / `admin123`
- **Staff**: `staff1` / `staff123`
- **Customer**: `customer1` / `customer123`
