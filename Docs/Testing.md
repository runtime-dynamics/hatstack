# Testing Guide

This document describes the testing strategy and practices for the H.A.T. Stack Bootstrap.

**📖 Main Documentation:** See the [main README](../README.md) for project overview and getting started.

**🐳 Deployment:** This project includes a production-ready `Dockerfile` for deployment to Google Cloud Run or any container platform.

## Table of Contents

- [Overview](#overview)
- [Running Tests](#running-tests)
- [Test Structure](#test-structure)
- [Writing Tests](#writing-tests)
- [Test Coverage](#test-coverage)
- [Continuous Integration](#continuous-integration)

## Overview

The bootstrap includes comprehensive tests for all major components:

- **Unit Tests** - Test individual functions and methods
- **Integration Tests** - Test component interactions (marked with `-short` flag)
- **Benchmark Tests** - Measure performance
- **Table-Driven Tests** - Test multiple scenarios efficiently

## Running Tests

### Quick Start

**Run all tests:**
```bash
# Linux/Mac
make test
# or
./test.sh test

# Windows
test.bat test
```

### Test Commands

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make test-verbose` | Run with verbose output |
| `make test-coverage` | Generate coverage report |
| `make test-race` | Run with race detector |
| `make test-short` | Skip integration tests |
| `make bench` | Run benchmarks |
| `make clean` | Clean test cache |

### Platform-Specific Commands

**Linux/Mac:**
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run benchmarks
make bench

# Using test.sh script
./test.sh test
./test.sh coverage
./test.sh race
```

**Windows:**
```cmd
REM Run all tests
test.bat test

REM Run with coverage
test.bat coverage

REM Run benchmarks
test.bat bench
```

### Running Specific Tests

```bash
# Test a specific package
go test ./config

# Test a specific function
go test ./config -run TestLoadConfig

# Test with verbose output
go test -v ./...

# Test with coverage for specific package
go test -cover ./config
```

## Test Structure

### Directory Layout

```
.
├── config/
│   ├── config.go
│   └── config_test.go
├── services/
│   ├── service.go
│   └── service_test.go
├── data/
│   ├── data.go
│   └── data_test.go
├── web/
│   ├── api/
│   │   ├── routes.go
│   │   └── routes_test.go
│   └── app/
│       ├── home.go
│       └── home_test.go
└── testutil/
    ├── testutil.go
    └── testutil_test.go
```

### Test File Naming

- Test files end with `_test.go`
- Test functions start with `Test`
- Benchmark functions start with `Benchmark`
- Example functions start with `Example`

## Writing Tests

### Basic Test Structure

```go
func TestFunctionName(t *testing.T) {
    // Arrange - Set up test data
    input := "test"
    expected := "result"
    
    // Act - Execute the function
    result := FunctionName(input)
    
    // Assert - Verify the result
    if result != expected {
        t.Errorf("FunctionName(%v) = %v, want %v", input, result, expected)
    }
}
```

### Table-Driven Tests

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case 1", "input1", "output1"},
        {"case 2", "input2", "output2"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Function(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Using Testify

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestWithTestify(t *testing.T) {
    result := Function()
    
    assert.NotNil(t, result)
    assert.Equal(t, expected, result)
    assert.Contains(t, result, "substring")
}
```

### Testing HTTP Handlers

```go
func TestHandler(t *testing.T) {
    // Create test router
    router := gin.New()
    router.GET("/test", YourHandler)
    
    // Create test request
    req, _ := http.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    
    // Execute request
    router.ServeHTTP(w, req)
    
    // Assert response
    assert.Equal(t, http.StatusOK, w.Code)
    assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}
```

### Mocking External Dependencies

For tests that require external services (like Datastore), use the `-short` flag to skip:

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Integration test code here
}
```

Run unit tests only:
```bash
go test -short ./...
```

### Benchmark Tests

```go
func BenchmarkFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Function()
    }
}

func BenchmarkFunctionWithSetup(b *testing.B) {
    // Setup (not measured)
    data := setupTestData()
    
    b.ResetTimer() // Reset timer after setup
    
    for i := 0; i < b.N; i++ {
        Function(data)
    }
}
```

Run benchmarks:
```bash
go test -bench=. -benchmem ./...
```

## Test Coverage

### Generating Coverage Reports

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

### Using Make

```bash
make test-coverage
# Opens coverage.html in your browser
```

### Coverage Goals

- **Overall**: Aim for 80%+ coverage
- **Critical paths**: 90%+ coverage
- **Config/Utils**: 95%+ coverage
- **Handlers**: 70%+ coverage (UI tests are harder)

### Viewing Coverage

The HTML report shows:
- **Green** - Covered lines
- **Red** - Uncovered lines
- **Gray** - Not executable

## Test Utilities

The `testutil` package provides helpers:

```go
import "runtime-dynamics/testutil"

// Create test router
router := testutil.SetupTestRouter()

// Create test context
w, c := testutil.CreateTestContext()

// Assert helpers
testutil.AssertNoError(t, err, "operation failed")
testutil.AssertError(t, err, "expected error")
```

## Best Practices

### DO

✅ Write tests for all public functions
✅ Use table-driven tests for multiple scenarios
✅ Test error cases and edge cases
✅ Use meaningful test names
✅ Keep tests independent and isolated
✅ Use `t.Helper()` in test utilities
✅ Clean up resources in `defer` statements
✅ Use `t.Parallel()` for independent tests

### DON'T

❌ Test private functions directly
❌ Make tests dependent on each other
❌ Use sleep for timing (use channels/contexts)
❌ Ignore test failures
❌ Write tests without assertions
❌ Test implementation details
❌ Commit commented-out tests

## Continuous Integration

Tests run automatically on:
- Every push to main branch
- Every pull request
- Multiple platforms (Windows, Linux, macOS)
- Multiple Go versions (1.21, 1.22)

See `.github/workflows/ci.yml` for CI configuration.

### CI Test Commands

```yaml
# Generate templates
templ generate

# Run tests
go test -v ./...

# Run with race detector
go test -race ./...

# Check coverage
go test -coverprofile=coverage.out ./...
```

## Troubleshooting

### Tests Fail Locally But Pass in CI

- Check Go version: `go version`
- Clean test cache: `go clean -testcache`
- Regenerate templates: `templ generate`
- Check environment variables

### Race Detector Warnings

```bash
# Run with race detector
go test -race ./...

# Fix data races by:
# - Using mutexes for shared data
# - Using channels for communication
# - Avoiding shared mutable state
```

### Coverage Not Updating

```bash
# Clean and regenerate
go clean -testcache
rm coverage.out coverage.html
make test-coverage
```

### Slow Tests

```bash
# Identify slow tests
go test -v ./... | grep -E "PASS|FAIL"

# Run only fast tests
go test -short ./...

# Profile tests
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

## Examples

### Complete Test Example

```go
package mypackage

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCompleteExample(t *testing.T) {
    // Table-driven test
    tests := []struct {
        name    string
        input   int
        want    int
        wantErr bool
    }{
        {"positive", 5, 25, false},
        {"zero", 0, 0, false},
        {"negative", -5, 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Square(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.want, result)
        })
    }
}

func BenchmarkSquare(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Square(42)
    }
}
```

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Go Test Coverage](https://go.dev/blog/cover)

## Contributing Tests

When contributing, ensure:

1. All new code has tests
2. Tests pass locally: `make test`
3. Coverage doesn't decrease: `make test-coverage`
4. No race conditions: `make test-race`
5. Tests are documented
6. Follow existing test patterns

See [CONTRIBUTING.md](../CONTRIBUTING.md) for more details.
