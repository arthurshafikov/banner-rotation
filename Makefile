BIN := "./.bin/app"
DOCKER_COMPOSE_FILE := "./deployments/docker-compose.yml"
DOCKER_COMPOSE_TEST_FILE := "./deployments/docker-compose.tests.yml"
APP_NAME := "banner_rotation"
APP_TEST_NAME := "banner_rotation_test"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -a -o $(BIN) -ldflags "$(LDFLAGS)" cmd/app/main.go

run: build 
	 $(BIN)

test: 
	go test --short -race ./internal/... ./pkg/...

lint: 
	golangci-lint run ./...

.PHONY: build test

enterdb:
	docker exec -it ${APP_NAME}_db_1 psql -U homestead;

up:
	docker-compose --env-file ./.env -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} up --build --attach app

down:
	docker-compose --env-file ./.env -f ${DOCKER_COMPOSE_FILE} down --volumes

mocks:
	mockgen -source=./internal/repository/repository.go -destination ./internal/repository/mocks/mock.go
	mockgen -source=./internal/services/services.go -destination ./internal/services/mocks/mock.go
	mockgen -source=./internal/transport/http/handler/handler.go -destination ./internal/transport/http/handler/mocks/mock.go
	mockgen -source=./internal/transport/http/server.go -destination ./internal/transport/http/mocks/mock.go
	mockgen -source=./pkg/queue/kafka.go -destination ./pkg/queue/mocks/mock.go

integration-tests:
	docker-compose --env-file ./.env.ci -f ${DOCKER_COMPOSE_TEST_FILE} -p ${APP_TEST_NAME} up  --build --abort-on-container-exit --exit-code-from integration --attach integration
	docker-compose -f ${DOCKER_COMPOSE_TEST_FILE} -p ${APP_TEST_NAME} down --volumes

reset-integration-tests:
	docker-compose -f ${DOCKER_COMPOSE_TEST_FILE} -p ${APP_TEST_NAME} down --volumes
