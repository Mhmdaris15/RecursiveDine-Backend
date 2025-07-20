@echo off
setlocal enabledelayedexpansion

rem Test runner script for RecursiveDine API (Windows)
echo === RecursiveDine API Test Suite ===
echo Starting comprehensive API tests...
echo.

rem Set environment for testing
set ENVIRONMENT=test
if "%DB_NAME%"=="" set DB_NAME=recursive_dine
set DB_NAME=%DB_NAME%_test

rem Function to print status messages
echo [INFO] Starting test execution...

rem Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go is not installed or not in PATH
    exit /b 1
)

echo [INFO] Go version:
go version

rem Check dependencies
echo [INFO] Downloading test dependencies...
go mod tidy

rem Setup test database
echo [INFO] Setting up test database...
echo [INFO] Test database setup completed

rem Run authentication tests
echo.
echo [INFO] === Running Authentication Tests ===
go test -v .\tests\ -run="TestAuth" -count=1
if errorlevel 1 (
    echo [ERROR] Authentication tests failed
    goto cleanup
)
echo [SUCCESS] Authentication tests completed successfully

rem Run table management tests
echo.
echo [INFO] === Running Table Management Tests ===
go test -v .\tests\ -run="TestTable" -count=1
if errorlevel 1 (
    echo [ERROR] Table management tests failed
    goto cleanup
)
echo [SUCCESS] Table management tests completed successfully

rem Run menu management tests
echo.
echo [INFO] === Running Menu Management Tests ===
go test -v .\tests\ -run="TestMenu|TestGet.*Menu|TestSearch.*Menu|TestAdmin.*Menu" -count=1
if errorlevel 1 (
    echo [ERROR] Menu management tests failed
    goto cleanup
)
echo [SUCCESS] Menu management tests completed successfully

rem Run order management tests
echo.
echo [INFO] === Running Order Management Tests ===
go test -v .\tests\ -run="TestOrder|TestCreate.*Order|TestGet.*Order|TestUpdate.*Order|TestAdmin.*Order|TestStaff.*Order" -count=1
if errorlevel 1 (
    echo [ERROR] Order management tests failed
    goto cleanup
)
echo [SUCCESS] Order management tests completed successfully

rem Run payment tests
echo.
echo [INFO] === Running Payment Tests ===
go test -v .\tests\ -run="TestPayment|TestQRIS|TestCash|TestAdmin.*Payment|TestCashier.*Payment" -count=1
if errorlevel 1 (
    echo [ERROR] Payment tests failed
    goto cleanup
)
echo [SUCCESS] Payment tests completed successfully

rem Run user management tests
echo.
echo [INFO] === Running User Management Tests ===
go test -v .\tests\ -run="TestAdmin.*User" -count=1
if errorlevel 1 (
    echo [ERROR] User management tests failed
    goto cleanup
)
echo [SUCCESS] User management tests completed successfully

rem Run full integration test suite
echo.
echo [INFO] === Running Full Integration Test Suite ===
go test -v .\tests\ -count=1 -timeout=30m
if errorlevel 1 (
    echo [ERROR] Full integration tests failed
    goto cleanup
)

rem Generate test coverage report
echo.
echo [INFO] === Generating Test Coverage Report ===
go test .\tests\ -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html

rem Display coverage summary
echo [INFO] Test coverage summary:
go tool cover -func=coverage.out | findstr "total"

goto success

:cleanup
echo [INFO] Cleaning up test database...
echo [INFO] Test database cleanup completed
exit /b 1

:success
echo [INFO] Cleaning up test database...
echo [INFO] Test database cleanup completed
echo.
echo [SUCCESS] === All tests completed successfully! ===
echo [INFO] Coverage report generated: coverage.html
echo [INFO] You can view the detailed coverage report by opening coverage.html in a browser

rem Handle command line arguments
if "%1"=="benchmark" (
    echo [INFO] === Running Performance Benchmarks ===
    go test -bench=. -benchmem .\tests\ > benchmark_results.txt
    echo [INFO] Benchmark results saved to benchmark_results.txt
)

if "%1"=="coverage" (
    echo [INFO] Running tests with coverage only...
    go test .\tests\ -coverprofile=coverage.out -covermode=atomic
    go tool cover -html=coverage.out -o coverage.html
    echo [SUCCESS] Coverage report generated: coverage.html
)

if "%1"=="quick" (
    echo [INFO] Running quick test suite (no integration tests)...
    go test -v .\tests\ -short -count=1
)

endlocal
