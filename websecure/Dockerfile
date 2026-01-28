# Multi-stage build
FROM golang:1.25-alpine AS builder

# ARGs for Git metadata (passed from Makefile/CI)
ARG GIT_COMMIT=unknown
ARG GIT_BRANCH=unknown
ARG BUILD_DATE
ARG VERSION=latest

WORKDIR /app

# First copy only module files for dependency caching
COPY go.mod ./
# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the application
# CGO_ENABLED=0 ensures a static binary for distroless/alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Add build metadata to the builder stage
LABEL stage=builder \
      git.commit=$GIT_COMMIT \
      git.branch=$GIT_BRANCH \
      build.date=$BUILD_DATE

# Final stage using distroless for a smaller attack surface
FROM gcr.io/distroless/static-debian11

# ARGs for metadata (required in this stage as well to populate labels)
ARG GIT_COMMIT=unknown
ARG GIT_BRANCH=unknown
ARG BUILD_DATE
ARG VERSION=latest

# OCI standard labels for traceability
LABEL org.opencontainers.image.title="Go Secure Server" \
      org.opencontainers.image.description="Secure Go web server for K8s security lab" \
      org.opencontainers.image.url="https://github.com/cooler-sai/GoK8sSecurityLab" \
      org.opencontainers.image.source="https://github.com/cooler-sai/GoK8sSecurityLab" \
      org.opencontainers.image.revision="${GIT_COMMIT}" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.authors="cooler-sai" \
      maintainer="cooler-sai"

# Optional: Custom labels for convenience
LABEL com.github.cooler-sai.project="GoK8sSecurityLab" \
      com.github.cooler-sai.component="go-secure-server" \
      com.github.cooler-sai.git-branch="${GIT_BRANCH}"

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/main .

# Run as non-privileged user for security
USER nonroot:nonroot

# Health check (recommended for Kubernetes/Docker monitoring)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/main", "health"]

EXPOSE 8080

# Use exec form to ensure signals (like SIGTERM) are handled correctly
ENTRYPOINT ["./main"]