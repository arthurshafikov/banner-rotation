BIN := "./.bin/app"
DOCKER_COMPOSE_FILE := "./deployments/docker-compose.yml"
DOCKER_COMPOSE_TEST_FILE := "./deployments/docker-compose.tests.yml"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -a -o $(BIN) -ldflags "$(LDFLAGS)" cmd/main.go

run: build 
	 $(BIN)

test: 
	go test --short -race ./internal/... ./pkg/...

.PHONY: build test

enterdb:
	docker exec -it banner-rotation_db_1 psql -U homestead;

migrate:
	goose -dir migrations postgres "host=localhost user=homestead password=secret dbname=homestead sslmode=disable" up

rollback:
	goose -dir migrations postgres "host=localhost user=homestead password=secret dbname=homestead sslmode=disable" down

fresh:
	goose -dir migrations postgres "host=localhost user=homestead password=secret dbname=homestead sslmode=disable" reset
	make migrate

mocks:
	mockgen -source=./internal/repository/repository.go -destination ./internal/repository/mocks/mock.go
	mockgen -source=./internal/services/services.go -destination ./internal/services/mocks/mock.go
