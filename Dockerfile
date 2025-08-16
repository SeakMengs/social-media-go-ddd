ARG GO_VERSION=1.24.6
FROM golang:${GO_VERSION} 
WORKDIR /app

RUN apt-get update && \
    apt-get install -y vim

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download