# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum from project root
COPY go.mod go.sum ./
# Download dependencies first to leverage Docker cache
RUN go mod download

# Copy only the necessary source code for the link-service and its internal dependencies
COPY internal/ ./internal/
COPY services/link-service/ ./services/link-service/

# Build a static, production-ready binary
# The source path must now be relative to the WORKDIR
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o link-service ./services/link-service

# Final stage
FROM scratch
WORKDIR /app
COPY --from=builder /app/link-service .
EXPOSE 8001
CMD ["./link-service"]
