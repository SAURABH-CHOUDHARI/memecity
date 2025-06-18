# Stage 1: Build the application
FROM golang:alpine AS builder

# Set working directory
WORKDIR /app

# Copy only the files needed for go mod
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app ./cmd/main.go

# Stage 2: Create a minimal production image
FROM alpine:3.18

# Add CA certificates and tzdata for security/timezone support
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose the port your application runs on
EXPOSE 8080

# Command to run the application
CMD ["./app"]