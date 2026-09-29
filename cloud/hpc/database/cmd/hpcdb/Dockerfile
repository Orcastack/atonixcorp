# -----------------------------
# Stage 1: Build the Go binary
# -----------------------------
FROM golang:1.27.1 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire HPC database module
COPY . .

# Build the hpcdb binary
RUN cd hpc/database/cmd/hpcdb && \
    go build -o /hpcdb

# -----------------------------
# Stage 2: Minimal runtime image
# -----------------------------
FROM debian:stable-slim

# Install CA certificates (needed for HTTPS)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the compiled binary
COPY --from=builder /hpcdb /usr/local/bin/hpcdb

# Expose API port
EXPOSE 8080

# Default command: run the API server
CMD ["hpcdb", "serve"]
