FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o sentinel main.go

FROM alpine:3.18
WORKDIR /
COPY --from=builder /app/sentinel /sentinel
# future folder for webhook certs
RUN mkdir -p /etc/webhook/certs
EXPOSE 8443
ENTRYPOINT ["/sentinel"]