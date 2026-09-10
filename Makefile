DB_URL ?= postgres://postgres:postgres@localhost:5432/subscriptions?sslmode=disable

.PHONY: run test lint tidy up down migrate-up migrate-down

run:
	go run .

test:
	go test ./... -cover

lint:
	golangci-lint run

tidy:
	go mod tidy

compose-up:
	docker compose up --build -d

compose-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1
