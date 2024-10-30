# Use the official Golang image as a parent image
FROM golang:1.22.6-alpine

# Set the working directory inside the container
WORKDIR /app

# Install air for hot reloading and postgresql-client
RUN apk add --no-cache postgresql-client && \
    go install github.com/cosmtrek/air@v1.49.0

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Add the startup script
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Expose port 8080 to the outside world
EXPOSE 8080

# Use air for development
CMD ["air", "-c", ".air.toml"]
