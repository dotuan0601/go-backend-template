run-dev:
	go run ./cmd/api -env dev

run-prod:
	go run ./cmd/api -env prod

test:
	go test ./...

build:
	mkdir -p bin
	go build -o bin/api ./cmd/api

migrate-up:
	go run ./cmd/migrate -direction up


migrate-down:
	go run ./cmd/migrate -direction down
