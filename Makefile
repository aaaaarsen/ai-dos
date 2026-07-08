.PHONY: run build docker docker-down migrate test

run:
	go run cmd/api/main.go

build:
	go build -o bin/server ./cmd/api

docker:
	docker compose up --build

docker-down:
	docker compose down

migrate:
	migrate -path migrations -database "postgres://postgres:9091@localhost:5432/ai_dos_dev?sslmode=disable" -verbose up

migrate-down:
	migrate -path migrations -database "postgres://postgres:9091@localhost:5432/ai_dos_dev?sslmode=disable" -verbose down

test:
	go test ./...

lint:
	go vet ./...