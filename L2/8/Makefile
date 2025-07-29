.PHONY: all build run vet lint install-lint

# Default target runs static analysis
all: vet lint

# Build the application binary
build:
	@go build -o myapp .

# Run the application
run: build
	@./myapp

# Run go vet to check for programmatic errors
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run golint to check for style issues
# Assumes golint is in ~/go/bin/golint
lint: install-lint
	@echo "Running golint..."
	@if ! command -v ~/go/bin/golint &> /dev/null; then \
		echo "golint not found. Please run 'make install-lint'"; \
		exit 1; \
	fi
	@~/go/bin/golint ./...

# Install the golint tool
install-lint:
	@echo "Installing golint..."
	@go install golang.org/x/lint/golint@latest
