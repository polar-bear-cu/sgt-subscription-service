.PHONY: run build test lint tidy up down

run:
	go run .

build:
	go build -o bin/server .

test:
	go test ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

up:
	docker compose up -d

down:
	docker compose down
