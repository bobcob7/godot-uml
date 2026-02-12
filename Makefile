BINARY := godot-uml
GOBIN  := $(shell pwd)/bin

.PHONY: all build test lint fmt serve clean generate install-tools

all: lint test build

build:
	go build -o $(GOBIN)/$(BINARY) ./cmd/godot-uml

test:
	go test -race -parallel 8 ./...

lint: fmt
	GOBIN=$(GOBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint
	$(GOBIN)/golangci-lint run ./...

fmt:
	GOBIN=$(GOBIN) go install mvdan.cc/gofumpt
	$(GOBIN)/gofumpt -w .

serve: build
	$(GOBIN)/$(BINARY) serve --port 8080

clean:
	rm -rf $(GOBIN)
	go clean ./...

generate:
	go generate ./...

install-tools:
	GOBIN=$(GOBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint
	GOBIN=$(GOBIN) go install mvdan.cc/gofumpt
