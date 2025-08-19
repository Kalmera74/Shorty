# Dockerfile for Shorty

FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Install git for module fetching (optional if not using private repos)
RUN apk add --no-cache git

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the Go app
RUN go build -o shorty ./cmd/server

# Expose service port (update if needed)
EXPOSE 8080

# Run the binary
CMD ["./shorty"]

