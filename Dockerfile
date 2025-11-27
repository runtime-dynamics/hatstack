# Use the official Golang image to create a build artifact.
# This is known as a multi-stage build.
FROM mirror.gcr.io/library/golang:latest AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install templ CLI for generating templates
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy the entire project
COPY . .

# Generate templ templates
RUN templ generate

# Build the Go app
# CGO_ENABLED=0 is required for a static build
# -o /app/server builds the application into a binary named "server" in the /app directory
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o /app/server ./cmd

# ---

# Start from a new, minimal image
FROM gcr.io/distroless/static-debian11
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/server /app/server

# Copy static files
COPY --from=builder /app/static /app/static

# Expose port 8080 to the outside world
EXPOSE 8080

# Set environment variable for Cloud Run compatibility
# Cloud Run sets PORT, but app uses LISTEN_PORT
ENV LISTEN_PORT=8080

# Command to run the executable
CMD ["/app/server"]
