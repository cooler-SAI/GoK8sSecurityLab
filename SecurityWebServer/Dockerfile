# Multi-stage build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# First copy only module files for dependency caching
COPY go.mod ./
# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage with distroless
FROM gcr.io/distroless/static-debian11

WORKDIR /app
COPY --from=builder /app/main .

USER nonroot:nonroot

EXPOSE 8080
CMD ["./main"]