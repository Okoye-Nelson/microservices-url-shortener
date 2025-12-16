# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
# Download dependencies first to leverage Docker cache
RUN go mod download

# Copy only the necessary source code
COPY internal/ ./internal/
COPY services/stats-service/ ./services/stats-service/

# Build a static, production-ready binary
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /stats-service ./services/stats-service/

# Final stage
FROM scratch

# Create a non-root user
USER 65532:65532

COPY --from=builder /stats-service /stats-service
EXPOSE 8003
CMD ["/stats-service"]
