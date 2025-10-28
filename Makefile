# Makefile for Deep Database Comparator

# Variables
BINARY_NAME=deepComparator
BUILD_DIR=./cmd
OUTPUT_DIR=.

# Build the application
build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(BUILD_DIR)

# Clean build artifacts
clean:
	rm -f $(OUTPUT_DIR)/$(BINARY_NAME)

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Lint code
vet:
	go vet ./...

# Run all checks
check: fmt vet test

# Build and run with example table
run-example: build
	./$(BINARY_NAME) -table=billing_model -verbose

# Create .env from example if it doesn't exist
setup:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from .env.example"; \
		echo "Please edit .env with your database configurations"; \
	else \
		echo ".env file already exists"; \
	fi

# Full setup: dependencies, build, and create .env
install: deps setup build
	@echo "Installation complete!"
	@echo "Next steps:"
	@echo "1. Edit .env with your database configurations"
	@echo "2. Run: make run-example"

# Development cycle
dev: clean check build

# Show help
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  clean       - Remove build artifacts"
	@echo "  deps        - Install Go dependencies"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format Go code"
	@echo "  vet         - Run Go vet"
	@echo "  check       - Run fmt, vet, and test"
	@echo "  run-example - Build and run with example table"
	@echo "  setup       - Create .env from example"
	@echo "  install     - Full setup (deps, setup, build)"
	@echo "  dev         - Development cycle (clean, check, build)"
	@echo "  help        - Show this help"

.PHONY: build clean deps test fmt vet check run-example setup install dev help