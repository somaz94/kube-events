BINARY_NAME=kube-events
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/somaz94/kube-events/cmd/cli.version=$(VERSION) -X github.com/somaz94/kube-events/cmd/cli.commit=$(COMMIT) -X github.com/somaz94/kube-events/cmd/cli.date=$(DATE)"

.PHONY: build clean test test-unit cover lint fmt vet demo demo-clean demo-all help

build: ## Build the binary
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/main.go

clean: ## Remove build artifacts
	rm -f $(BINARY_NAME) coverage.out

test: test-unit ## Run all tests (alias for test-unit)

test-unit: ## Run unit tests with coverage
	go test ./... -v -race -cover -coverprofile=coverage.out

cover: test-unit ## Generate coverage report
	go tool cover -func=coverage.out

cover-html: test-unit ## Open coverage report in browser
	go tool cover -html=coverage.out

lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

demo: build ## Run demo (deploy test resources + show events)
	./scripts/demo.sh

demo-clean: ## Clean up demo resources
	./scripts/demo-clean.sh

demo-all: demo demo-clean ## Run demo and clean up

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
