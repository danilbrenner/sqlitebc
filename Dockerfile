FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

RUN apk add --no-cache git && apk add build-base && apk add gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .. .

RUN export CGO_ENABLED=1 && go build -o sqlitebc .

FROM --platform=$BUILDPLATFORM alpine:latest

WORKDIR /app

COPY --from=builder /app/sqlitebc .
