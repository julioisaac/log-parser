SRC=$(shell find . -name "*.go")
IMAGE_NAME=log-parser
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

check-dependencies:
	go mod tidy
	git diff --exit-code go.mod

deps: check-dependencies
	go mod tidy
	go mod download
	go install -v github.com/go-critic/go-critic/cmd/gocritic@latest

test:
	go test -v -cover `go list ./...`

run:
	go run main.go

build:
	docker build -t log-parser .

start:
	docker run --name log-parser --env-file .env log-parser

stop:
	docker stop log-parser

lint:
	gocritic check ./...
