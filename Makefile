.DEFAULT_GOAL := help

COMMIT_SHA := $(shell git describe --dirty --always --tags --long)

RED := $(shell tput setaf 1)
GREEN := $(shell tput setaf 2)
YELLOW := $(shell tput setaf 3)
CYAN := $(shell tput setaf 6)
NC := $(shell tput sgr0)

COVERAGE_FILE := /tmp/coverage.out

STATIC_CHECK := go run honnef.co/go/tools/cmd/staticcheck@latest
VULN_CHECK := go run golang.org/x/vuln/cmd/govulncheck@latest
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest

##@ Misc

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: confirm
confirm: ## Ask for confirmation before continuing
	@echo '$(YELLOW)Are you sure?$(NC) [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty: ## Check if there are uncommitted changes
	@git diff --quiet --exit-code || (echo "$(RED)Uncommitted changes$(NC)" && exit 1)

##@ Code quality

.PHONY: format
format: ## Format code and tidy modfile
	@go fmt ./...
	@go mod tidy -v

.PHONY: lint
lint: ## Run linters
	@go vet ./...
	@$(STATIC_CHECK) -checks=all,-ST1000,-U1000 ./...
	@$(GOLANGCI_LINT) run ./...

.PHONY: audit
audit: format lint ## Run all code quality checks
	@go mod verify
	@$(VULN_CHECK) ./...

##@ Development

.PHONY: test
test: ## Run all tests
	@go test -v -race -buildvcs ./...

.PHONY: test-cover
test-cover: ## Run all tests and display coverage
	@go test -v -race -buildvcs -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE)

##@ Operations

tag ?= $(shell bash -c 'read -p "Tag: " tag; echo $$tag')
.PHONY: release
release: no-dirty confirm ## Create a new release
	@git tag -a $(tag)
	@git push --tags
