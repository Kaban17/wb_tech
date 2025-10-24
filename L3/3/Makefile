.PHONY: all build run vet lint install-lint

# Default target runs static analysis
all: vet lint

# Build the application binary
build:
	@go build -o myapp ./cmd/app/main.go

# Run the application
run: build
	@./myapp

test:
	@go test -v ./...

# Run go vet to check for programmatic errors
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run golint to check for style issues
# Assumes golint is in ~/go/bin/golint
lint:
	@echo "Running golint..."
	golangci-lint run
