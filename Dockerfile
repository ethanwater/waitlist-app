# Build stage
FROM golang:1.21.1 AS build

WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the remaining files
COPY . .

# Set environment variables for cross-compilation to Linux
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build the Go binary
RUN go build -o waitlist

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/waitlist .

# Copy the static files (HTML, CSS, etc.)
COPY --from=build /app/views /app/views

# Ensure the binary has execute permissions
RUN chmod +x waitlist

# Expose the port your application runs on (if needed)
EXPOSE 8080

# Run the binary
CMD ["./waitlist"]

