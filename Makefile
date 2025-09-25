# Ministo Go Project Makefile
# Build automation for Go cryptocurrency miner

# Variables
GO		= go
BIN		= ministo
SRC		= ./cmd/ministo
BUILD_DIR	= ./bin
LDFLAGS		= -s -w

# Default target - show help
.DEFAULT_GOAL := help

.PHONY: help all build clean install test fmt vet lint run deps mod-tidy mod-download

help:	## Show this help
	@echo "Ministo - Go Cryptocurrency Miner"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

all: clean fmt vet build	## Clean, format, vet and build

build:	## Build the project
	@echo "Building $(BIN)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BIN) $(SRC)

run: build	## Build and run the project
	@echo "Running $(BIN)..."
	@$(BUILD_DIR)/$(BIN)

install:	## Install binary to GOPATH/bin
	$(GO) install -ldflags="$(LDFLAGS)" $(SRC)

test:	## Run tests
	$(GO) test -v ./...

fmt:	## Format Go code
	$(GO) fmt ./...

vet:	## Run go vet
	$(GO) vet ./...

lint:	## Run golangci-lint (requires installation)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

deps:	## Download dependencies
	$(GO) mod download

mod-tidy:	## Tidy and verify dependencies
	$(GO) mod tidy
	$(GO) mod verify

mod-download:	## Download dependencies
	$(GO) mod download

clean:	## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@$(GO) clean -cache -testcache -modcache 2>/dev/null || true

info:	## Show project information
	@echo "Project: Ministo Go Cryptocurrency Miner"
	@echo "Binary: $(BIN)"
	@echo "Source: $(SRC)"
	@echo "Build: $(BUILD_DIR)"
	@echo "Go version: $$($(GO) version)"