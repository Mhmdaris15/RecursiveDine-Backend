# Testing Documentation

This document provides comprehensive information about the testing setup for the RecursiveDine API.

## Overview

The test suite includes:
- **Unit Tests**: Individual component testing
- **Integration Tests**: End-to-end API endpoint testing
- **Authentication Tests**: JWT and role-based access control
- **Database Tests**: Repository and service layer testing
- **Coverage Reports**: Code coverage analysis

## Test Structure

```
tests/
├── api_test_setup.go      # Test suite setup and utilities
├── auth_test.go           # Authentication endpoint tests
├── table_test.go          # Table management tests
├── menu_test.go           # Menu management tests
├── order_test.go          # Order management tests
├── payment_test.go        # Payment processing tests
└── user_test.go           # User management tests
```

## Running Tests

### Prerequisites

1. **Go 1.21+** installed
2. **PostgreSQL** database running
3. **Environment variables** set:
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=your_user
   export DB_PASSWORD=your_password
   export DB_NAME=recursive_dine
   export JWT_SECRET=your_jwt_secret
   ```

### Quick Start

**Windows:**
```powershell
.\run_tests.bat
```

**Linux/macOS:**
```bash
chmod +x run_tests.sh
./run_tests.sh
```

### Test Commands

#### Run All Tests
```bash
go test ./tests/ -v
```

#### Run Specific Test Categories
```bash
# Authentication tests
go test ./tests/ -run="TestAuth" -v

# Table management tests
go test ./tests/ -run="TestTable" -v

# Menu management tests
go test ./tests/ -run="TestMenu" -v

# Order management tests
go test ./tests/ -run="TestOrder" -v

# Payment tests
go test ./tests/ -run="TestPayment" -v

# User management tests
go test ./tests/ -run="TestUser" -v
```

#### Generate Coverage Report
```bash
go test ./tests/ -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html
```

#### Run Benchmarks
```bash
go test -bench=. -benchmem ./tests/
```

### Test Runner Options

The test runner scripts support several options:

**Coverage Only:**
```bash
./run_tests.sh coverage    # Linux/macOS
.\run_tests.bat coverage   # Windows
```

**Quick Tests (no integration):**
```bash
./run_tests.sh quick       # Linux/macOS
.\run_tests.bat quick      # Windows
```

**Benchmarks:**
```bash
./run_tests.sh benchmark   # Linux/macOS
.\run_tests.bat benchmark  # Windows
```

## Test Categories

### 1. Authentication Tests (`auth_test.go`)

Tests JWT-based authentication system:

- **User Registration**: Valid/invalid registration attempts
- **User Login**: Credential validation and token generation
- **Token Refresh**: Refresh token functionality
- **Logout**: Token invalidation
- **Input Validation**: Email format, password strength

**Key Test Cases:**
```go
TestAuthRegister()
TestAuthLogin() 
TestAuthRefreshToken()
TestAuthLogout()
TestAuthLoginInvalidCredentials()
```

### 2. Table Management Tests (`table_test.go`)

Tests restaurant table management:

- **QR Code Access**: Public table lookup by QR code
- **Admin Operations**: CRUD operations for tables
- **Status Management**: Table availability status updates
- **Authorization**: Role-based access control

**Key Test Cases:**
```go
TestTableGetByQRCode()
TestAdminGetAllTables()
TestAdminCreateTable()
TestAdminUpdateTable()
TestAdminUpdateTableStatus()
TestAdminDeleteTable()
```

### 3. Menu Management Tests (`menu_test.go`)

Tests menu and category management:

- **Public Access**: Menu viewing for customers
- **Search Functionality**: Item search with filters
- **Category Management**: CRUD operations for categories
- **Item Management**: CRUD operations for menu items
- **Availability Control**: Staff can update item availability

**Key Test Cases:**
```go
TestGetMenu()
TestSearchMenuItems()
TestAdminCreateCategory()
TestAdminCreateMenuItem()
TestUpdateMenuItemAvailability()
```

### 4. Order Management Tests (`order_test.go`)

Tests order processing system:

- **Order Creation**: Customer order placement
- **Order Retrieval**: User and admin order access
- **Status Updates**: Order workflow management
- **Admin Features**: Order statistics and management
- **Authorization**: Role-based order access

**Key Test Cases:**
```go
TestCreateOrder()
TestGetUserOrders()
TestAdminGetAllOrders()
TestUpdateOrderStatus()
TestGetOrderStatistics()
```

### 5. Payment Tests (`payment_test.go`)

Tests payment processing:

- **QRIS Payments**: Digital payment initiation and verification
- **Cash Payments**: Cashier cash transaction processing
- **Payment Status**: Payment tracking and status updates
- **Refunds**: Admin refund processing
- **Reconciliation**: Cashier shift reconciliation

**Key Test Cases:**
```go
TestInitiateQRISPayment()
TestVerifyPayment()
TestCashierProcessCashPayment()
TestAdminProcessRefund()
TestCashierReconcilePayments()
```

### 6. User Management Tests (`user_test.go`)

Tests user administration:

- **User CRUD**: Admin user management operations
- **Role Management**: User role assignments
- **Access Control**: Authorization testing
- **Data Validation**: Input validation and constraints

**Key Test Cases:**
```go
TestAdminGetAllUsers()
TestAdminCreateUser()
TestAdminUpdateUser()
TestAdminDeleteUser()
```

## Test Data Setup

The test suite automatically creates test users with different roles:

- **Admin**: `admin@test.com` - Full system access
- **Staff**: `staff@test.com` - Order and menu management
- **Cashier**: `cashier@test.com` - Payment processing
- **Customer**: `customer@test.com` - Basic ordering

Each test user has a corresponding JWT token for authenticated requests.

## Database Testing

### Test Database

Tests use a separate database (`{DB_NAME}_test`) to avoid affecting production data.

### Migration and Cleanup

- **Setup**: Automatic table creation via GORM auto-migration
- **Cleanup**: Tables are dropped after test completion
- **Isolation**: Each test suite runs in a clean database state

## Coverage Goals

Target coverage levels:
- **Overall**: >80%
- **Controllers**: >85%
- **Services**: >90%
- **Critical Paths**: 100% (authentication, payments)

## Continuous Integration

### GitHub Actions (Example)

```yaml
name: Test Suite
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_DB: recursive_dine_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: '1.21'
    - name: Run tests
      run: ./run_tests.sh
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

## Performance Testing

### Benchmark Tests

Performance benchmarks test:
- **Response Times**: API endpoint latency
- **Memory Usage**: Memory allocation patterns
- **Throughput**: Requests per second capacity

### Load Testing

For load testing, use tools like:
- **Apache Benchmark (ab)**
- **wrk**
- **Vegeta**

Example load test:
```bash
ab -n 1000 -c 10 -H "Authorization: Bearer $TOKEN" \
   http://localhost:8002/api/v1/menu
```

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   - Verify PostgreSQL is running
   - Check environment variables
   - Ensure test database exists

2. **Authentication Failures**
   - Check JWT_SECRET environment variable
   - Verify token generation in test setup

3. **Permission Errors**
   - Confirm role-based middleware is working
   - Check test user creation and token assignment

4. **Port Conflicts**
   - Ensure test server port is available
   - Check for running instances

### Debug Mode

Enable verbose output:
```bash
go test ./tests/ -v -args -debug
```

### Test Database Inspection

Connect to test database:
```bash
psql -h localhost -U your_user -d recursive_dine_test
```

## Contributing

### Adding New Tests

1. **Create test file**: Follow naming convention `*_test.go`
2. **Add to test suite**: Include in `api_test_setup.go`
3. **Follow patterns**: Use existing test patterns and utilities
4. **Update documentation**: Add test descriptions here

### Test Best Practices

1. **Isolation**: Each test should be independent
2. **Cleanup**: Clean up test data after each test
3. **Assertions**: Use descriptive assertion messages
4. **Coverage**: Aim for comprehensive test coverage
5. **Performance**: Keep tests fast and efficient

### Code Review Checklist

- [ ] Tests cover happy path and error cases
- [ ] Authorization tests included for protected endpoints
- [ ] Input validation tests present
- [ ] Database cleanup implemented
- [ ] Test documentation updated

## Results Interpretation

### Coverage Report

The HTML coverage report shows:
- **Green**: Well-covered code (>80%)
- **Yellow**: Moderately covered code (50-80%)
- **Red**: Poorly covered code (<50%)
- **Gray**: Uncovered code

### Test Results

Success indicators:
- All tests pass (green checkmarks)
- Coverage above target thresholds
- No memory leaks in benchmarks
- Reasonable response times

Failure indicators:
- Test failures (red X marks)
- Low coverage percentages
- High memory allocation
- Slow response times

Focus on red and yellow areas for improvement.
