.PHONY: all build test bench clean install uninstall fmt lint coverage release help

# Variables
BINARY_NAME := fast-cc-hooks
GC_BINARY_NAME := gc
BUILD_DIR := build
CMD_DIR := cmd/fast-cc-hooks
GC_CMD_DIR := cmd/gc
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GOFLAGS := -v
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.commit=$(COMMIT) -w -s"

# Go tools versions
GOLANGCI_LINT_VERSION := v1.55.2
GORELEASER_VERSION := latest

# Default target
all: clean lint test build

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^##' Makefile | sed -E 's/^## /  /'

## build: Build the binary for current platform
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## build-gc: Build the gc helper utility
build-gc:
	@echo "Building $(GC_BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(GC_BINARY_NAME) ./$(GC_CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(GC_BINARY_NAME)"

## build-all-tools: Build all tools
build-all-tools: build build-gc

## build-all: Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux AMD64
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	# Linux ARM64
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)
	# macOS AMD64
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	# macOS ARM64 (M1/M2)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	# Windows AMD64
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "Multi-platform build complete"

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v -race -timeout 30s ./...

## test-short: Run short tests
test-short:
	@echo "Running short tests..."
	@go test -short -v ./...

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem -run=^$$ ./...

## coverage: Generate test coverage report
coverage:
	@echo "Generating coverage report..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

## lint: Run linters
lint:
	@echo "Running linters..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	fi
	@golangci-lint run --timeout 5m

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## mod: Download and tidy modules
mod:
	@echo "Downloading and tidying modules..."
	@go mod download
	@go mod tidy
	@go mod verify

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## install: Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME) to GOPATH/bin..."
	@go install $(LDFLAGS) ./$(CMD_DIR)
	@echo "Installation complete"

## install-gc: Install the gc utility to GOPATH/bin
install-gc: build-gc
	@echo "Installing $(GC_BINARY_NAME) to GOPATH/bin..."
	@go install $(LDFLAGS) ./$(GC_CMD_DIR)
	@echo "Installation complete"

## install-all: Install all tools
install-all: install install-gc

## uninstall: Remove the binary from GOPATH/bin
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(GOPATH)/bin/$(BINARY_NAME)
	@rm -f $(GOPATH)/bin/$(GC_BINARY_NAME)
	@echo "Uninstall complete"

## run: Build and run the binary
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

## dev: Run with hot reload (requires air)
dev:
	@if ! command -v air &> /dev/null; then \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@air

## release: Create release with goreleaser (requires goreleaser)
release:
	@if ! command -v goreleaser &> /dev/null; then \
		echo "Installing goreleaser..."; \
		go install github.com/goreleaser/goreleaser@$(GORELEASER_VERSION); \
	fi
	@goreleaser release --clean

## release-snapshot: Create snapshot release
release-snapshot:
	@if ! command -v goreleaser &> /dev/null; then \
		echo "Installing goreleaser..."; \
		go install github.com/goreleaser/goreleaser@$(GORELEASER_VERSION); \
	fi
	@goreleaser release --snapshot --clean

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) -t $(BINARY_NAME):latest .

## ci: Run CI pipeline locally
ci: clean mod fmt vet lint test build
	@echo "CI pipeline complete"

## check: Quick check before commit
check: fmt vet test-short
	@echo "Pre-commit checks complete"

# Installation targets for git hooks
## hook-install: Install pre-commit hook
hook-install: build
	@./$(BUILD_DIR)/$(BINARY_NAME) install

## hook-uninstall: Uninstall pre-commit hook  
hook-uninstall: build
	@./$(BUILD_DIR)/$(BINARY_NAME) uninstall

## init-config: Initialize configuration file
init-config: build
	@./$(BUILD_DIR)/$(BINARY_NAME) init