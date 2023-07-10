SHELL=/bin/bash
GOPACKAGES=$(shell go list ./... | egrep -v "vendor|mock")

install: check-dependencies vendor tools

check-dependencies:
	go mod tidy
	git diff --exit-code go.mod
	git diff --exit-code go.sum

vendor:
	go mod vendor

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

checks: test check-dependencies lint

test:
	go test -v -cover `go list ./...`

lint:
	golangci-lint run -v

deps: check-dependencies
	go mod tidy
	go mod download

run:
	go run main.go

docker/build:
	docker build -t log-parser .

docker/start: docker/stop
	docker run --name log-parser log-parser

docker/stop:
	docker stop log-parser
	docker rm log-parser
