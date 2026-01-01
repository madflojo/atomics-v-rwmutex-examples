.PHONY: all clean tests lint build format benchmarks

all: build tests lint

# Run tests with coverage
tests:
	@echo "Running tests with coverage..."
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@if command -v go tool cover > /dev/null 2>&1; then \
		go tool cover -html=coverage.out -o coverage.html; \
	fi

# Run benchmarks
benchmarks:
	@echo "Running benchmarks..."
	go test -run=^$$ -bench=. -benchmem ./...

# Build the package
build:
	@echo "Building package..."
	go build ./...

# Format code
format:
	@echo "Formatting code..."
	gofmt -s -w .
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not installed, skipping goimports"; \
	fi
	@if command -v golines >/dev/null 2>&1; then \
		golines -w .; \
	else \
		echo "golines not installed, skipping golines"; \
	fi

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping lint"; \
	fi
