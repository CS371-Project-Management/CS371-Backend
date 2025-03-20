FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

WORKDIR /app/cmd

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates mysql-client

WORKDIR /app

COPY --from=builder /app/cmd/main /app/main

COPY .env .env

EXPOSE 8080

CMD ["/app/main"]
