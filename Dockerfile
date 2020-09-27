# Start from golang base image
FROM golang:1.14 as builder

ENV GO111MODULE=on

# Add Maintainer info
LABEL Company="NOVA Solutions Co"
LABEL maintainer="Ahmed Abouzied <ahmedaabouzied44@gmail.com>"
LABEL Project="Tasarruf mobile application backend"

# Install git.
# Git is required for fetching the dependencies.
# RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /tasarruf

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tasarruf .

# Start a new stage from scratch
FROM alpine:3.10.4
# RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /tasarruf/tasarruf .
COPY --from=builder /tasarruf/.env .

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
ENTRYPOINT ["./tasarruf"]
