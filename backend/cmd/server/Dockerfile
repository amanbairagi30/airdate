# Use the official Golang image as a parent image
FROM golang:1.22.6-alpine

# Set the working directory inside the container
WORKDIR /app

# Add necessary packages
RUN apk add --no-cache postgresql-client

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Add a startup script
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/start.sh"]
