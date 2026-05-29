# ============================================
# Stage 1: Build the Go binary
# ============================================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install dependencies for CGO (if needed)
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
# CGO_ENABLED=0 creates a fully static binary
# -ldflags="-w -s" strips debug info for smaller size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	-ldflags="-w -s" \
	-o main ./cmd/api/main.go

# ============================================
# Stage 2: Create minimal runtime image
# ============================================
FROM alpine:3.21

WORKDIR /app

# Install CA certificates for HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
adduser -u 1001 -S appuser -G appgroup

# Default runtime port, can be overridden by docker-compose or `-e PORT=...`
ENV PORT=8080

# Copy binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check (optional but recommended)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
	CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

	# Run the binary
CMD ["./main"]