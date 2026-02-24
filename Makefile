.PHONY: build test lint fmt docker

build:
	go build ./...

test:
	go test ./...

lint:
	go vet ./...

fmt:
	gofmt -w $(shell find . -name '*.go')

docker:
	docker build -t zenx:latest .
