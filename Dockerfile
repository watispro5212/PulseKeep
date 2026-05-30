# ==========================================================
# STAGE 1: Builder
# ==========================================================
FROM golang:1.25-alpine AS builder

# Install build essentials
RUN apk add --no-cache git ca-certificates build-base

WORKDIR /app

# Copy dependency manifests and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build statically compiled binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/pulsekeep \
    ./cmd/pulsekeep

# ==========================================================
# STAGE 2: Runner
# ==========================================================
FROM alpine:3.19 AS runner

# Install essential runtime utilities (certificates for SSL, timezone DB)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/pulsekeep /pulsekeep

# Expose default HTTP Port
EXPOSE 7860
ENV PORT=7860

# Run the app
ENTRYPOINT ["/pulsekeep"]
