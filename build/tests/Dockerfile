FROM golang:1.17.3-alpine as build

RUN set -ex; \
    apk update; \
    apk add --no-cache git; \
    apk add build-base

ENV GOOSE_FILE /bin/goose
ENV MIGRATIONS /migrations

COPY ./migrations ${MIGRATIONS}

ADD https://github.com/pressly/goose/releases/download/v3.3.0/goose_linux_x86_64 ${GOOSE_FILE}
RUN chmod +x ${GOOSE_FILE}

ENV CODE_DIR /go/src/
ENV TEST_DIR ./tests

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait /wait
RUN chmod +x /wait

CMD /wait && \
    ${GOOSE_FILE} -dir ${MIGRATIONS} postgres "${DSN}" reset && \
    ${GOOSE_FILE} -dir ${MIGRATIONS} postgres "${DSN}" up && \
    go test -v ${TEST_DIR}
