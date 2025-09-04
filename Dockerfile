# Base stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Install swag and generate docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --dir ./cmd/api,./internal --output ./docs

# Build API binary
RUN go build -o shorty-api ./cmd/api

# Build worker binary
RUN go build -o analytics-worker ./cmd/analytics-worker

# -----------------------------
# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/shorty-api .
COPY --from=builder /app/analytics-worker .
COPY --from=builder /app/docs ./docs

# Expose API port
EXPOSE 8080

# Default command: run API
CMD ["./shorty-api"]
