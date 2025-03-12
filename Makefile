# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	
	
	@go build .

# Run the application
run:
	@go run .

# Test the application
test: lint
	@echo "Testing..."
	@go test ./... -v

# Test coverage
coverage:
	@echo "Coverage..."
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# Linting and security checks
lint:
	@echo "Running linting checks..."
	@staticcheck ./...
	@echo "Running security checks..."
	@gosec ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f pdfmc


.PHONY: all build run test coverage linting clean
