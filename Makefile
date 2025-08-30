# Makefile for Shorty

APP_NAME=shorty
CMD_PATH=cmd/shorty
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

.PHONY: build run test fmt tidy vet docker docker-run clean api all

build:
	go build -o bin/$(APP_NAME) ./$(CMD_PATH)

run:
	go run ./$(CMD_PATH)

test:
	go test ./...

fmt:
	go fmt $(GO_FILES)

vet:
	go vet ./...

tidy:
	go mod tidy

docker:
	docker compose up

api:
	swag init --dir ./cmd/shorty,./internal --output ./docs

all: tidy fmt vet build test

clean:
	rm -rf bin/
