FROM golang:1.24.1-alpine
WORKDIR /app
RUN apk add --no-cache gcc musl-dev mysql-client

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]