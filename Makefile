.PHONY: build test lint fmt clean install

VERSION ?= dev
LDFLAGS := -s -w -X main.version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/jira ./cmd/jira
	go build -ldflags "$(LDFLAGS)" -o bin/confluence ./cmd/confluence

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

clean:
	rm -rf bin/

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/jira
	go install -ldflags "$(LDFLAGS)" ./cmd/confluence
