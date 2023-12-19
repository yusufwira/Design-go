# Use an official Go runtime as a base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download and install any dependencies
RUN go mod download

# Build the Go application
RUN go build -o main .

# Expose the port the application uses
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

