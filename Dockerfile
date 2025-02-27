FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grok-proxy ./cmd/server

# Use a small image for the final container
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/grok-proxy .
# Copy config files
COPY --from=builder /app/configs ./configs

# Set environment variables
ENV GIN_MODE=release

EXPOSE 8000

# Run the application
CMD ["./grok-proxy"]
