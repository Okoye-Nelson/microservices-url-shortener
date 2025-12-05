# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum from project root
COPY go.mod go.sum ./
RUN go mod download

# Copy service-specific code
COPY services/link-service/*.go ./
RUN go build -o link-service

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/link-service .
EXPOSE 8001
CMD ["./link-service"]
