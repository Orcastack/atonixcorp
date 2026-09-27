
FROM golang:1.27.1 AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY cloud/compute/ ./compute/


COPY cloud/internal/ ./internal/


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o compute-service ./compute/main.go



FROM debian:13-slim

LABEL org.opencontainers.image.source="https://github.com/orcastack/atonixcorp"

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/compute-service /app/compute-service

EXPOSE 8580

CMD ["/app/compute-service"]
