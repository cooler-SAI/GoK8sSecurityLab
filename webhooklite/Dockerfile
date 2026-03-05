FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -o webhooklite main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/webhooklite .
EXPOSE 8080
CMD ["./webhooklite"]