FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .. .

RUN go build -o sqlitebc .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/sqlitebc .
