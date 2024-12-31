# Use the official Golang image as a builder
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules manifests first (to leverage Docker caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application's source code
COPY . .

# Build the Go application
RUN go build -o app ./cmd/main.go

# Use a minimal image for the final output
FROM debian:bullseye-slim

# Set working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app .

# Expose the port the app runs on (replace 8080 with your application's port if different)
EXPOSE 8080

# Run the application
CMD ["./app"]
