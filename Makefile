BIN := "./.bin/app"
DOCKER_COMPOSE_FILE := "./deployments/docker-compose.yml"
DOCKER_COMPOSE_TEST_FILE := "./deployments/docker-compose.tests.yml"
APP_NAME := "banner_rotation"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -a -o $(BIN) -ldflags "$(LDFLAGS)" cmd/app/main.go

run: build 
	 $(BIN)

test: 
	go test --short -race ./internal/... ./pkg/...

.PHONY: build test

enterdb:
	docker exec -it ${APP_NAME}_db_1 psql -U homestead;

up:
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} up --build

down:
	docker-compose -f ${DOCKER_COMPOSE_FILE} down --volumes

mocks:
	mockgen -source=./internal/repository/repository.go -destination ./internal/repository/mocks/mock.go
	mockgen -source=./internal/services/services.go -destination ./internal/services/mocks/mock.go
	mockgen -source=./internal/transport/http/handler/handler.go -destination ./internal/transport/http/handler/mocks/mock.go
