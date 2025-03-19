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

# Environment variables for configuration
# - PBI_SERVER_PORT: Port to listen on (default: 3000)
# - PBI_SRV_API_KEY: API key for authentication
# - MINIO_ENDPOINT: MinIO server endpoint
# - MINIO_ACCESS_KEY: MinIO access key
# - MINIO_SECRET_KEY: MinIO secret key
# - MINIO_FRAGMENTS_BUCKET: MinIO bucket for fragment files (default: ifc-fragment-files)
# - MINIO_LCA_COST_DATA_BUCKET: MinIO bucket for LCA cost data (default: lca-cost-data)
ENV PBI_SERVER_PORT=3000

# Expose port to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./main"]