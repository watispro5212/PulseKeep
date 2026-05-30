FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/pulsekeep \
    ./cmd/pulsekeep

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/pulsekeep /pulsekeep

EXPOSE 7860
ENV PORT=7860

ENTRYPOINT ["/pulsekeep"]
