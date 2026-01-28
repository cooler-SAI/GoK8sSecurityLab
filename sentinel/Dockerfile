# Use specific Alpine version for reproducible builds
FROM golang:1.25-alpine3.23 AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download and verify dependencies
RUN go mod download && go mod verify

# Copy all source files (not just main.go)
COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-w -s -extldflags '-static'" \
    -trimpath \
    -o sentinel .

# --- Final image ---
FROM alpine:3.23.3

# Create non-root user for security
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

WORKDIR /

# Copy binary with correct permissions
COPY --from=builder --chown=appuser:appuser /app/sentinel /sentinel

# Create directories with proper permissions
RUN mkdir -p /etc/webhook/certs && \
    chown -R appuser:appuser /etc/webhook && \
    chmod 755 /sentinel

# Switch to non-root user
USER appuser

# Health check for container orchestration
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8443/health || exit 1

# Container metadata labels (VERSION will be passed during build)
LABEL org.opencontainers.image.source="https://github.com/your/repo" \
      org.opencontainers.image.description="Secure webhook sentinel"

# Expose webhook port
EXPOSE 8443

# Use exec form for proper signal handling
ENTRYPOINT ["/sentinel"]

# Optional: Default command arguments
# CMD ["--config", "/config.yaml"]