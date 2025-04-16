# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@go build .

# Run the application
install: build
	@echo "Installing..."
	@sudo mv pdfmc /usr/local/bin/
	@mkdir -p ~/.zsh/completions
	@pdfmc completion zsh > ~/.zsh/completions/_pdfmc

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
	@if [ -f /usr/local/bin/pdfmc ]; then
		sudo rm /usr/local/bin/pdfmc
	fi

gif: install
	@echo "Recreating the Gifs..."
	@cd public && vhs merge.tape
	@cd public && vhs encrypt.tape
	@cd public && vhs completions.tape
	@cd public && vhs decrypt.tape


.PHONY: all build install run test coverage linting clean gif
