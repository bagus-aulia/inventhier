# Stage 1: Build
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api/main.go

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies (for database connectivity)
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/api .

# Expose port (matches your SERVER_PORT)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./api"]
