# Use official Go image
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . ./

# Build the application
RUN go build -o /app/bin/api ./cmd/api

# Final stage - minimal Alpine image
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/bin/api /root/api

# Expose port (Railway will set PORT env var)
EXPOSE 8080

# Run the binary
CMD ["/root/api"]
