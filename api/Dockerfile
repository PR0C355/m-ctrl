# syntax=docker/dockerfile:1

FROM golang:1.22.5

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go .

ENV GIN_MODE=release

RUN go build -o /app/api

CMD ["/app/api"]