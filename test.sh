#!/bin/bash

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

function show_help {
    echo "Available commands:"
    echo "  ./test.sh test      - Run all tests"
    echo "  ./test.sh verbose   - Run tests with verbose output"
    echo "  ./test.sh coverage  - Run tests with coverage report"
    echo "  ./test.sh race      - Run tests with race detector"
    echo "  ./test.sh short     - Run only short tests"
    echo "  ./test.sh bench     - Run benchmarks"
    echo "  ./test.sh clean     - Clean test cache"
}

function run_tests {
    echo -e "${GREEN}Running tests...${NC}"
    go test ./...
}

function run_verbose {
    echo -e "${GREEN}Running tests (verbose)...${NC}"
    go test -v ./...
}

function run_coverage {
    echo -e "${GREEN}Running tests with coverage...${NC}"
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}Coverage report generated: coverage.html${NC}"
}

function run_race {
    echo -e "${GREEN}Running tests with race detector...${NC}"
    go test -race ./...
}

function run_short {
    echo -e "${GREEN}Running short tests...${NC}"
    go test -short ./...
}

function run_bench {
    echo -e "${GREEN}Running benchmarks...${NC}"
    go test -bench=. -benchmem ./...
}

function clean {
    echo -e "${GREEN}Cleaning test cache...${NC}"
    go clean -testcache
    rm -f coverage.out coverage.html
}

case "$1" in
    test)
        run_tests
        ;;
    verbose)
        run_verbose
        ;;
    coverage)
        run_coverage
        ;;
    race)
        run_race
        ;;
    short)
        run_short
        ;;
    bench)
        run_bench
        ;;
    clean)
        clean
        ;;
    *)
        show_help
        ;;
esac
