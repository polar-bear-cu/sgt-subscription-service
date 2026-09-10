DB_URL ?= postgres://postgres:postgres@localhost:5433/subscriptions?sslmode=disable

.PHONY: run test lint tidy up down migrate-up migrate-down

run:
	swag init -g main.go -o docs --parseInternal
	go run .

test:
	go test ./... -cover

format:
	golangci-lint fmt

lint:
	golangci-lint run

tidy:
	go mod tidy

docs:
	swag init -g main.go -o docs --parseInternal

compose-up:
	docker compose up --build -d --wait

compose-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-create:
	@test -n "$(name)" || (echo "usage: make migrate-create name=<snake>"; exit 1)
	migrate create -ext sql -dir migrations "$(name)"