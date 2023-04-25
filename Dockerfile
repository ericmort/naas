# Use the official Golang image as the base image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the dependencies using Go modules
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Golang application
RUN go build -o main .

# Expose the port the application will run on
EXPOSE 8080

# Run the compiled binary
CMD ["./main"]
