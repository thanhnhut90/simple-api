# Use the official Golang image as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files into the container
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the application code into the container
COPY . .

# install build stuff
RUN apk update && apk upgrade
RUN apk add --no-cache make build-base

# Build the Go application
RUN export CGO_ENABLED=1; go build ./cmd/api

# Use a smaller base image for the final image
FROM alpine:latest

# Copy the binary from the previous stage
COPY --from=0 /app/api /app/api

RUN apk update && apk upgrade
RUN apk add --no-cache sqlite

# Expose the port the application will listen on
EXPOSE 8080

# Run the application when the container starts
ENTRYPOINT ["/app/api"]
CMD ["-db=sqlite"]

