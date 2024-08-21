# Stage 1: Build the Go binary
FROM golang:alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Run the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the build stage
COPY --from=build /app/main .

# Expose the port the app will run on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
