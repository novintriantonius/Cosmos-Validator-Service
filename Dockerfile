# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o cosmos-validator-service cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/cosmos-validator-service .

# Create a non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./cosmos-validator-service"] 