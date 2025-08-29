#!/bin/bash

# Проверка аргументов
if [ $# -ne 2 ]; then
  echo "Usage: $0 <base_dir_name> <N>"
  echo "Example: $0 L42 5"
  exit 1
fi

BASEDIR="$1"
N="$2"

# Проверка, что N — число
if ! [[ "$N" =~ ^[0-9]+$ ]]; then
  echo "Error: N must be a positive integer."
  exit 2
fi

for ((i = 1; i <= N; i++)); do
  PROJECT_DIR="$BASEDIR/$i"
  mkdir -p "$PROJECT_DIR"
  cd "$PROJECT_DIR" || exit 3

  # Инициализация Go-модуля
  go mod init "wb_tech/${BASEDIR,,}_$i" > /dev/null 2>&1

  # Создание main.go
  cat > main.go <<EOF
package main

import "fmt"

func main() {
    fmt.Println("Hello, world from ${BASEDIR}/${i}!")
}
EOF

  # Создание Makefile (ВНУТРИ директории проекта)
  cat > Makefile <<EOF
.PHONY: all build run vet lint install-lint

# Default target runs static analysis
all: vet lint

# Build the application binary
build:
	@go build -o myapp .

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
EOF

  cd - > /dev/null
done

echo "✅ Created $N Go projects in $BASEDIR/{1..$N}"
