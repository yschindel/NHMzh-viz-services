# Use the official Golang image as the base image
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o main .

# Final stage
FROM ubuntu:latest

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
# Copy SQL files from builder stage
COPY --from=builder /app/storage/azure/sql/ /app/storage/azure/sql/

ENV PBI_SERVER_PORT=3000

# Expose port to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./main"]