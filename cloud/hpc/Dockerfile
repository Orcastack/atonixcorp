# ------------------------------------------------------------
# 1. Build Stage (Debian 13 + Go 1.27.1)
# ------------------------------------------------------------
FROM golang:1.27.1 AS builder

# Install Go 1.27.1 and build tools
RUN apt-get update && apt-get install -y wget git build-essential ca-certificates && \
    wget https://go.dev/dl/go1.27.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.27.1.linux-amd64.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

# Copy module files first (cache optimization)
COPY go.mod go.sum ./
RUN go mod download

# Copy entire project
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o hpc ./cmd/hpc

# ------------------------------------------------------------
# 2. Runtime Stage (Debian Slim)
# ------------------------------------------------------------
FROM debian:13-slim

RUN apt-get update && apt-get install -y ca-certificates && \
    useradd -m hpcuser

WORKDIR /opt/hpc

# Copy binary
COPY --from=builder /app/hpc /opt/hpc/hpc

# Copy config
COPY config/hpc.yaml /opt/hpc/config/hpc.yaml

# Permissions
RUN chown -R hpcuser:hpcuser /opt/hpc

USER hpcuser

EXPOSE 8080 9090

ENTRYPOINT ["/opt/hpc/hpc"]
