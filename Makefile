BINARY_NAME=kube-events
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/somaz94/kube-events/cmd/cli.version=$(VERSION) -X github.com/somaz94/kube-events/cmd/cli.commit=$(COMMIT) -X github.com/somaz94/kube-events/cmd/cli.date=$(DATE)"

.PHONY: build clean test test-unit cover lint fmt vet demo demo-clean demo-all check-gh branch pr help

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

## Workflow

check-gh: ## Check if gh CLI is installed and authenticated
	@command -v gh >/dev/null 2>&1 || { echo "\033[31m✗ gh CLI not installed. Run: brew install gh\033[0m"; exit 1; }
	@gh auth status >/dev/null 2>&1 || { echo "\033[31m✗ gh CLI not authenticated. Run: gh auth login\033[0m"; exit 1; }
	@echo "\033[32m✓ gh CLI ready\033[0m"

branch: ## Create feature branch (usage: make branch name=watch-mode)
	@if [ -z "$(name)" ]; then echo "Usage: make branch name=<feature-name>"; exit 1; fi
	git checkout main
	git pull origin main
	git checkout -b feat/$(name)
	@echo "\033[32m✓ Branch feat/$(name) created\033[0m"

pr: check-gh ## Run tests, push, and create PR (usage: make pr title="Add watch mode")
	@if [ -z "$(title)" ]; then echo "Usage: make pr title=\"PR title\""; exit 1; fi
	go test ./... -race -cover
	go vet ./...
	git push -u origin $$(git branch --show-current)
	@./scripts/create-pr.sh "$(title)"
	@echo "\033[32m✓ PR created\033[0m"

## Help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
