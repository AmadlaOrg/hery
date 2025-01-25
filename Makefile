.PHONY: install-deps
install-deps: ## Installs Dependencies
	@echo "--->  Installing Dependencies"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/boumenot/gocover-cobertura@latest
	@go install github.com/jstemmer/go-junit-report/v2@latest
	@go install github.com/jandelgado/gcov2lcov@latest
	@go install github.com/vektra/mockery/v3@latest

generate: ## Generate mock code
	@echo "--->  Generating code"
	@go generate ./...
	@go run github.com/vektra/mockery/v2@latest

.PHONY: lint
lint: ## Linting
	@echo "--->  Linting"
	@golangci-lint run -v

.PHONY: lint-fix
lint-fix: ## Lint-Fixing code
	@echo "---> Lint-Fixing code"
	@golangci-lint run --fix

.PHONY: test
test: ## Test code
	@.script/test.sh

.PHONY: cov
cov: cov ## Show test coverage
	@go tool cover -html=.reports/coverage.out

.PHONY: test-cov
test-cov: test cov ## Test coverage

.PHONY: clean
clean: ## Clean bin and coverage files
	@echo "--->  Cleaning bin and coverage files"
	@rm -f bin/*
	@rm -f coverage.out
	@rm -f .reports/*

build: ## Build code
	@echo "---> Build"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -buildvcs=true -o bin/ ./

.PHONY: help
help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sed 's/Makefile://' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
