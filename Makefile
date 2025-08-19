# Makefile for Shorty

APP_NAME=shorty
CMD_PATH=cmd/server
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

.PHONY: build run test lint fmt tidy docker docker-run clean

build:
	go build -o bin/$(APP_NAME) ./$(CMD_PATH)

run:
	go run ./$(CMD_PATH)

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

tidy:
	go mod tidy

docker:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

clean:
	rm -rf bin/

