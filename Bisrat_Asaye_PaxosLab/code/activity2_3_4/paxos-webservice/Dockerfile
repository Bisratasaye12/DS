# Step 1: Build the Go binary
FROM golang:1.19 AS builder

# Set the working directory in the container
WORKDIR /app

# Copy the rest of the application code
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o paxos-webservice webservice.go

# Step 2: Create a smaller image to run the Go binary
FROM alpine:latest

# Install dependencies for running the app (certificates for HTTPS)
RUN apk update
RUN apk --no-cache add ca-certificates

# Set the working directory for the app
WORKDIR /

# Copy the Go binary from the builder image
COPY --from=builder /app/paxos-webservice .

ENTRYPOINT ["/paxos-webservice"]

# Expose the application port (8080 in this case)
EXPOSE 8080


