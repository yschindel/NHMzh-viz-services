# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code and SQL files
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o main .

# Final stage
FROM ubuntu:latest

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
# Copy SQL files from builder stage
COPY --from=builder /app/sql/ /app/sql/

# Run the application
CMD ["./main"]