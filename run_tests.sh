#!/bin/bash

# Test runner script for RecursiveDine API
set -e

echo "=== RecursiveDine API Test Suite ==="
echo "Starting comprehensive API tests..."
echo ""

# Set environment for testing
export ENVIRONMENT=test
export DB_NAME="${DB_NAME}_test"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

print_status "Go version: $(go version)"

# Check if database is running
print_status "Checking database connection..."

# Run dependency check
print_status "Downloading test dependencies..."
go mod tidy

# Run tests with different verbosity levels
run_tests() {
    local test_type=$1
    local test_pattern=$2
    local description=$3
    
    print_status "Running $description..."
    
    if go test -v ./tests/ -run="$test_pattern" -count=1; then
        print_success "$description completed successfully"
    else
        print_error "$description failed"
        return 1
    fi
}

# Create test database if it doesn't exist
setup_test_db() {
    print_status "Setting up test database..."
    
    # This would typically create a test database
    # For now, we assume the test database setup is handled in the test code
    print_status "Test database setup completed"
}

# Clean up test database
cleanup_test_db() {
    print_status "Cleaning up test database..."
    # Add cleanup logic here if needed
    print_status "Test database cleanup completed"
}

# Main test execution
main() {
    print_status "Starting test execution..."
    
    # Setup
    setup_test_db
    
    # Run different test suites
    echo ""
    print_status "=== Running Authentication Tests ==="
    if ! run_tests "auth" "TestAuth" "Authentication tests"; then
        cleanup_test_db
        exit 1
    fi
    
    echo ""
    print_status "=== Running Table Management Tests ==="
    if ! run_tests "table" "TestTable" "Table management tests"; then
        cleanup_test_db
        exit 1
    fi
    
    echo ""
    print_status "=== Running Menu Management Tests ==="
    if ! run_tests "menu" "TestMenu|TestGet.*Menu|TestSearch.*Menu|TestAdmin.*Menu" "Menu management tests"; then
        cleanup_test_db
        exit 1
    fi
    
    echo ""
    print_status "=== Running Order Management Tests ==="
    if ! run_tests "order" "TestOrder|TestCreate.*Order|TestGet.*Order|TestUpdate.*Order|TestAdmin.*Order|TestStaff.*Order" "Order management tests"; then
        cleanup_test_db
        exit 1
    fi
    
    echo ""
    print_status "=== Running Payment Tests ==="
    if ! run_tests "payment" "TestPayment|TestQRIS|TestCash|TestAdmin.*Payment|TestCashier.*Payment" "Payment tests"; then
        cleanup_test_db
        exit 1
    fi
    
    echo ""
    print_status "=== Running User Management Tests ==="
    if ! run_tests "user" "TestAdmin.*User" "User management tests"; then
        cleanup_test_db
        exit 1
    fi
    
    # Run all tests together for integration testing
    echo ""
    print_status "=== Running Full Integration Test Suite ==="
    if ! go test -v ./tests/ -count=1 -timeout=30m; then
        print_error "Full integration tests failed"
        cleanup_test_db
        exit 1
    fi
    
    # Generate test coverage report
    echo ""
    print_status "=== Generating Test Coverage Report ==="
    go test ./tests/ -coverprofile=coverage.out -covermode=atomic
    go tool cover -html=coverage.out -o coverage.html
    
    # Display coverage summary
    print_status "Test coverage summary:"
    go tool cover -func=coverage.out | tail -1
    
    # Cleanup
    cleanup_test_db
    
    echo ""
    print_success "=== All tests completed successfully! ==="
    print_status "Coverage report generated: coverage.html"
    print_status "You can view the detailed coverage report by opening coverage.html in a browser"
}

# Run performance benchmarks (optional)
run_benchmarks() {
    print_status "=== Running Performance Benchmarks ==="
    go test -bench=. -benchmem ./tests/ > benchmark_results.txt
    print_status "Benchmark results saved to benchmark_results.txt"
}

# Handle command line arguments
case "${1:-}" in
    "benchmark")
        run_benchmarks
        ;;
    "coverage")
        print_status "Running tests with coverage only..."
        go test ./tests/ -coverprofile=coverage.out -covermode=atomic
        go tool cover -html=coverage.out -o coverage.html
        print_success "Coverage report generated: coverage.html"
        ;;
    "quick")
        print_status "Running quick test suite (no integration tests)..."
        go test -v ./tests/ -short -count=1
        ;;
    *)
        main
        ;;
esac
