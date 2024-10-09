# Build stage
FROM golang:1.21.1 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum for dependency management
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application (assuming main.go is in the root)
RUN go build -o myapp ./main.go

# Final stage - use a smaller image
FROM alpine:latest

# Set working directory in the final image
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .

# Expose the application port if necessary
EXPOSE 8080

# Command to run the application
CMD ["./myapp"]
