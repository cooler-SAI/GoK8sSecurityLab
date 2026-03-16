FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -o webhook ./cmd/webhooklite

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/webhook .
COPY certs/ /certs/
EXPOSE 8443
CMD ["./webhook"]