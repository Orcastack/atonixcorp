# ─────────────────────────────────────────────
# Stage 1 — Build with Go 1.27.1
# ─────────────────────────────────────────────
FROM golang:1.27.1 AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build all daemons
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/smtpd ./cmd/smtpd
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/imapd ./cmd/imapd
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/pop3d ./cmd/pop3d
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/apid ./cmd/apid
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/admind ./cmd/admind
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/webmail ./cmd/webmail

# ─────────────────────────────────────────────
# Stage 2 — Runtime image
# ─────────────────────────────────────────────
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

# Copy binaries
COPY --from=builder /app/bin /app/bin

# Copy configs + DKIM keys + webmail assets
COPY configs /app/configs
COPY webmail /app/webmail

# Default daemon (override per service)
CMD ["/app/bin/smtpd"]
