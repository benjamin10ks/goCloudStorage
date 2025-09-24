# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application (your main.go is in cmd/server/)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/

# Final stage
FROM alpine:latest

# Install ca-certificates and create necessary directories
RUN apk --no-cache add ca-certificates && \
    mkdir -p /app/web /app/storage/users

# Create non-root user
RUN addgroup -g 1001 appgroup && adduser -u 1001 -G appgroup -s /bin/sh -D appuser

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy web assets (templates, static files)
COPY --from=builder /app/web ./web

# Create storage directory with proper permissions
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Run the application
CMD ["./main"]
