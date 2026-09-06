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