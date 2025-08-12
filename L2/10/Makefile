.PHONY: all lint vet test

all: lint vet test

lint:
	@echo "Running golint..."
	@go install golang.org/x/lint/golint@latest
	@golint ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

test:
	@echo "Running go test..."
	@go test ./...
