# Build stage
FROM golang:1.23-alpine AS builder

# Install dependencies for CGO and SQLite
RUN apk add --no-cache git build-base sqlite-dev

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Enable CGO and build both binaries
ENV CGO_ENABLED=1
RUN go build -ldflags="-s -w" -o gomer-server ./cmd/gomer/main.go
RUN go build -ldflags="-s -w" -o gocli ./cmd/cli/main.go

# Final image
FROM alpine:3.20

# Install required libraries for SQLite support and correct timezone handling
RUN apk add --no-cache sqlite-libs tzdata

# Copy both binaries from the builder stage
COPY --from=builder /app/gomer-server /usr/local/bin/gomer-server
COPY --from=builder /app/gocli /usr/local/bin/gocli

# Set working directory and entrypoint
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/gomer-server"]
