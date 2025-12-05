# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/stats-service/*.go ./
RUN go build -o stats-service

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/stats-service .
EXPOSE 8003
CMD ["./stats-service"]
