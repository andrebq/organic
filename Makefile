.PHONY: gotidy watch default test-containers

default: test

start-test-containers:
	docker run --rm -d -p 6379:6379 --name "organic-redis" redis:alpine

stop-test-containers:
	docker kill "organic-redis"

monitor-redis:
	docker exec -ti "organic-redis" /bin/sh

test:
	go test ./...

build:
	go build ./...

gotidy: build
	go fmt ./...
	go mod tidy

watch:
	modd

build-shell:
	go build ./cmd/oshell
