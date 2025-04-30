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
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-server ./cmd/api-server

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/api-server .

# Expose the API port
EXPOSE 8082

# Command to run the application
CMD ["./api-server"] 