# Stage 1 -

# Build the Go binary
FROM golang:1.24.4-alpine AS builder

# Install tool needed to build Go App
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Binary for CrawlerX
RUN go build -o crawlerx ./cmd/server

# Stage 2 -

# Run lightweight container
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/crawlerx .

# Set Port
ENV PORT=8080

# Expose port
EXPOSE 8080

# Run the app
CMD ["./crawlerx"]