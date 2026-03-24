# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 ensures a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/ghost ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/bin/ghost .

# Copy environment variables file if needed (optional)
COPY .env .env

# Expose ports
# REST
EXPOSE 8080
# gRPC
EXPOSE 50051

# Run the executable
CMD ["./ghost"]
