@echo off
REM Windows batch file for running tests

if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="test" goto test
if "%1"=="verbose" goto verbose
if "%1"=="coverage" goto coverage
if "%1"=="race" goto race
if "%1"=="short" goto short
if "%1"=="bench" goto bench
if "%1"=="clean" goto clean
goto help

:help
echo Available commands:
echo   test.bat test      - Run all tests
echo   test.bat verbose   - Run tests with verbose output
echo   test.bat coverage  - Run tests with coverage report
echo   test.bat race      - Run tests with race detector
echo   test.bat short     - Run only short tests
echo   test.bat bench     - Run benchmarks
echo   test.bat clean     - Clean test cache
goto end

:test
echo Running tests...
go test ./...
goto end

:verbose
echo Running tests (verbose)...
go test -v ./...
goto end

:coverage
echo Running tests with coverage...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
echo Coverage report generated: coverage.html
goto end

:race
echo Running tests with race detector...
go test -race ./...
goto end

:short
echo Running short tests...
go test -short ./...
goto end

:bench
echo Running benchmarks...
go test -bench=. -benchmem ./...
goto end

:clean
echo Cleaning test cache...
go clean -testcache
if exist coverage.out del coverage.out
if exist coverage.html del coverage.html
goto end

:end
