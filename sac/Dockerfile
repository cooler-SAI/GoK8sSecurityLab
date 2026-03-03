FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY . .
RUN go mod tidy || go mod init webhook
RUN go build -o webhook main.go


FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/webhook .

# DO NOT COPY CERTIFICATES! They will come from Kubernetes secret
# Remove these lines:
# COPY cert.pem .
# COPY key.pem .

EXPOSE 8443
CMD ["./webhook"]