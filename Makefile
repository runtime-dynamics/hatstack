.PHONY: help test test-verbose test-coverage test-race test-short bench clean build run dev

# Default target
help:
	@echo "Available targets:"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race     - Run tests with race detector"
	@echo "  test-short    - Run only short tests (skip integration)"
	@echo "  bench         - Run benchmarks"
	@echo "  clean         - Clean build artifacts and test cache"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with Air (live reload)"
	@echo "  generate      - Generate templ templates"

# Run all tests
test:
	@echo "Running tests..."
	@go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	@go test -race ./...

# Run only short tests (skip integration tests)
test-short:
	@echo "Running short tests..."
	@go test -short ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Clean build artifacts and test cache
clean:
	@echo "Cleaning..."
	@go clean -testcache
	@rm -f coverage.out coverage.html
	@rm -rf tmp/

# Generate templ templates
generate:
	@echo "Generating templ templates..."
	@templ generate

# Build the application
build: generate
	@echo "Building..."
	@go build -o ./tmp/main ./cmd

# Run the application
run: build
	@echo "Running application..."
	@./tmp/main

# Run with Air (live reload)
dev:
	@echo "Starting development server with Air..."
	@air
