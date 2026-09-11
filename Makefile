DB_URL ?= postgres://postgres:postgres@localhost:5433/subscriptions?sslmode=disable
IMAGE ?= sgt-subscription-service
PORT ?= 8080

.PHONY: run test lint format tidy docs image container compose-up compose-down migrate-up migrate-down migrate-create

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

image:
	docker build -t $(IMAGE) .

container: image
	docker run --rm --env-file .env -p $(PORT):8080 $(IMAGE)

compose-up:
	docker compose up -d --wait --remove-orphans

compose-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-create:
	@test -n "$(name)" || (echo "usage: make migrate-create name=<snake>"; exit 1)
	migrate create -ext sql -dir migrations "$(name)"