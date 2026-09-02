GOPATH=$(shell go env GOPATH)
IMAGE_REGISTRY=dockerhub
IMAGE_NAMESPACE ?= inventhier
IMAGE_NAME ?= $(shell basename `pwd`)
CURRENT_PATH=$(shell pwd)
COMMIT_ID ?= $(shell git rev-parse --short HEAD)
GO111MODULE=on

PROTOCDIR = ./api/proto
protoc: $(PROTOCDIR)/*/*
	for file in $^ ; do \
		protoc -I $${file} --go_out=. --go-grpc_out=. $${file}/*.proto; \
	done