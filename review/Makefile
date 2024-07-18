.PHONY: build run test generate install-lint lint lint-ci

DOCKER_IMAGE_NAME=coupon_service

deps:
	go mod download

build:
	go build -o bin/coupon_service cmd/coupon_service/main.go

install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.49.0

lint:
	golangci-lint run ./...

lint-ci:
	golangci-lint run ./... --out-format=colored-line-number --timeout=5m

test: generate
	go test ./...

generate:
	go generate ./...

run: generate
	go run ./cmd/coupon_service/main.go

docker-compose-up:
	docker-compose up --build -d

docker-compose-down:
	docker-compose down
