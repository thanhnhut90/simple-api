SHELL := /bin/bash
GOPATH ?= $(shell go env GOPATH)

default: help

# Help target
help:
	@echo "Available targets:"
	@echo "  build: Build the project"
	@echo "  test: Run tests"
	@echo "  fmt: Run fmt"

build:
	go build ./cmd/api

test:
	go test ./...

fmt:
	go list -f '{{.Dir}}' ./... | xargs gofmt -w

.PHONY: test fmt 
