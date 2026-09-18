GOPATH=$(shell go env GOPATH)
IMAGE_REGISTRY=dockerhub
IMAGE_NAMESPACE ?= inventhier
IMAGE_NAME ?= $(shell basename `pwd`)
CURRENT_PATH=$(shell pwd)
COMMIT_ID ?= $(shell git rev-parse --short HEAD)
GO111MODULE=on

# Generate proto file
PROTOCDIR = ./api/proto ./internal/adapters/client/grpc
protoc:
	for dir in $(PROTOCDIR); do \
		if [ -d "$$dir" ]; then \
			find "$$dir" -name "*.proto" -type f -print0 | xargs -0 -I {} protoc -I "$$dir" --go_out=. --go-grpc_out=. {}; \
		fi \
	done
	
# Install mockery v3.6.1
install-mockery:
	go install github.com/vektra/mockery/v3@v3.6.1

# Generate all mocks using go generate
generate: install-mockery
	go generate ./...

# Generate all mocks using mockery config
mocks: install-mockery
	mockery

# Install golangci-lint
install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
lint: install-linter
	golangci-lint run --timeout=5m ./...

# Run unit tests with verbose output
test:
	go get github.com/newm4n/goornogo
	export GO111MODULE on; \
	go test ./... -cover -vet=all -v -short -covermode=count -coverprofile=coverage.out > test.txt
	goornogo -i coverage.out -c 70

# Run unit tests with coverage report
test-coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run unit tests with short flag (faster)
test-short:
	go test -v -short ./...

# Run benchmarks
bench:
	go test -bench=. -benchmem -run=^$ ./...

# Lint and test together
test-lint: lint test

# Full quality check: lint, test, and coverage
quality: lint test-coverage
	@echo "Quality check complete!"

# Help target
help:
	@echo "Available targets:"
	@echo "  make protoc              - Generate proto files"
	@echo "  make install-mockery     - Install mockery v3.6.1"
	@echo "  make generate            - Generate all mocks using go generate"
	@echo "  make mocks               - Generate mocks using mockery config"
	@echo "  make install-linter      - Install golangci-lint"
	@echo "  make lint                - Run golangci-lint on all packages"
	@echo "  make test                - Run all unit tests with coverage"
	@echo "  make test-coverage       - Generate HTML coverage report"
	@echo "  make test-short          - Run tests with short flag (faster)"
	@echo "  make bench               - Run benchmarks"
	@echo "  make test-lint           - Run lint and tests"
	@echo "  make quality             - Full quality check (lint, test, coverage)"
	@echo "  make help                - Show this help message"

.PHONY: protoc install-mockery generate mocks install-linter lint test test-coverage test-short bench test-lint quality help