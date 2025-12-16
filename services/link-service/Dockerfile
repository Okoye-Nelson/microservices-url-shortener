# --- Build Stage ---
# Use the official Golang image to create a build-time environment.
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files first. This is a key optimization.
# Docker will cache this layer and only re-run it if go.mod or go.sum change.
COPY go.mod go.sum ./
RUN go mod download

# Copy only the necessary source code for this service
COPY internal/ ./internal/
COPY services/link-service/ ./services/link-service/

# Build the application.
# CGO_ENABLED=0 creates a statically linked binary.
# -ldflags="-w -s" strips debug information, reducing binary size.
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /main ./services/link-service/

# --- Final Stage ---
# Use a minimal, non-root base image for the final application.
FROM scratch

# Create a non-root user to run the application
# Although 'scratch' has no shell, user/group IDs are still respected by the kernel.
# Using a high UID/GID is a good practice to avoid conflicts with host system users.
USER 65532:65532


# Copy only the compiled binary from the builder stage
COPY --from=builder /main /main

# Expose the port the application will run on
EXPOSE 8001

# Set the command to run the application
CMD ["/main"]
