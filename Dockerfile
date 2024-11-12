# Stage 1: Build the Go application
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

# Stage 2: Runtime with Docker CLI, kubectl, and kind
FROM docker:latest

# Install dependencies: kubectl and kind
RUN apk add --no-cache curl bash && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/ && \
    curl -Lo /usr/local/bin/kind "https://kind.sigs.k8s.io/dl/latest/kind-linux-amd64" && \
    chmod +x /usr/local/bin/kind

# Copy the built Go application from the builder stage
COPY --from=builder /app/main /app/main



# Set the working directory
WORKDIR /app

# Copy the assets directory into the container
COPY ./assets /app/assets

# Set the PROJECT_ROOT environment variable
ENV PROJECT_ROOT=/app

# Expose port 8080 for the Go server
EXPOSE 8080

# Command to run the Go server
# command that will keep the container running
# CMD ["tail", "-f", "/dev/null"]
CMD ["./main"]
