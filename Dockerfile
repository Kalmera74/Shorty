# Dockerfile for Shorty
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Install git for module fetching (optional if not using private repos)
RUN apk add --no-cache git

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code, including docs
COPY . .

# Ensure Swagger docs are included (if they are generated in cmd/shorty/docs)
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --dir ./cmd/api,./internal --output ./docs

# Build the Go app
RUN go build -o shorty ./cmd/api

# Expose service port (update if needed)
EXPOSE 8080

# Run the binary
CMD ["./shorty"]
