# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/redirect-service/*.go ./
RUN go build -o redirect-service

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/redirect-service .
EXPOSE 8002
CMD ["./redirect-service"]
