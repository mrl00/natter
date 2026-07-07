.PHONY: run build test lint swagger docker-up docker-down docker-build clean

run:
	go run ./cmd/server

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server ./cmd/server
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/healthcheck ./cmd/healthcheck

test:
	hurl --test tests/e2e/*.hurl

lint:
	go vet ./...

swagger:
	swag init -g cmd/server/main.go

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-build:
	docker build -f build/docker/Dockerfile -t natter .

clean:
	rm -rf bin/* tmp/*
