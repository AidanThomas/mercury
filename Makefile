help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
.PHONY: help

run: ## Run cli client locally
	$(call setenv,debug)
	@go run cmd/cli/main.go 127.0.0.1:1234 testuser
.PHONY: run

run2:
	$(call setenv,debug)
	@go run cmd/cli/main.go 127.0.0.1:1234 testuser2
.PHONY: run

server: ## Run server locally
	$(call setenv,debug)
	@go run cmd/server/main.go 1234
.PHONY: server

dev: ## Start server locally and monitor for changes
	@find . -name '*.go' | entr -rcs 'make server'
.PHONY: dev

lint: ## Lint code
	@golint ./... | grep -v unexported || true
	@go vet ./... 2>&1 || echo ''
.PHONY: lint

clean: ## Clean the project
	@rm -rf bin
	@go clean
.PHONY: clean

define setenv
	$(eval include $1.env)
	$(eval export)
endef
