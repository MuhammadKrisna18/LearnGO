.PHONY: run build docker-up docker-down

run:
	go run cmd/api/main.go

run-front:
	cd frontend && npm run dev -- --open

build:
	go build -o bin/api cmd/api/main.go

docker-up:
	docker compose up -d

docker-down:
	docker compose down
