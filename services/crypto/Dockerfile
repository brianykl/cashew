# Use an Alpine image as a base
FROM alpine:latest AS builder

# Install dependencies
RUN apk add --no-cache bash wget tar gcc musl-dev go

# Set the Go version you want to install
ENV GO_VERSION 1.21.6

# Download and install Go from source
RUN wget https://dl.google.com/go/go$GO_VERSION.src.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.src.tar.gz && \
    cd /usr/local/go/src && \
    ./make.bash

# Set environment variables
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory in the container
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the application files to the working directory in the container
COPY . .

# Build the Go app
RUN go build -o crypto .

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory inside the container for the new stage
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/crypto .

# Expose the port that the microservice listens on (adjust as needed)
EXPOSE 5002

# Command to run the executable
CMD ["./crypto"]
