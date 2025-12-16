# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
 
# Copy only the necessary source code
COPY internal/ ./internal/
COPY services/redirect-service/ ./services/redirect-service/

# Build a static, production-ready binary
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o redirect-service ./services/redirect-service/

# Final stage
FROM scratch
WORKDIR /app
COPY --from=builder /app/redirect-service .
EXPOSE 8002
CMD ["./redirect-service"]
