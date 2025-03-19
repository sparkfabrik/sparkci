.PHONY: build clean test lint vet fmt mod help

APP_NAME := sparkci
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
GO_FILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
GOBIN ?= $(shell go env GOPATH)/bin
SRC_DIRS := cmd pkg

help: ## Show help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build $(LDFLAGS) -o $(APP_NAME) ./cmd/sparkci

install: build ## Install the application
	cp $(APP_NAME) $(GOBIN)/

clean: ## Clean build artifacts
	rm -f $(APP_NAME)
	rm -rf ./dist

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linters
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed"; \
		exit 1; \
	fi

vet: ## Run go vet
	go vet ./...

fmt: ## Format source code
	gofmt -s -w $(GO_FILES)

mod: ## Tidy and verify dependencies
	go mod tidy
	go mod verify

dev: ## Run the application in development mode
	go run ./cmd/sparkci

ci: lint test build ## Run CI tasks

.DEFAULT_GOAL := help