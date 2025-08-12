# Use a multi-platform base image (ARM and AMD64)

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM
WORKDIR /app

# Copy the entire current directory into the container's working directory
COPY . .

# Build the Go application
RUN go build -o picture-api .


# Install curl for debugging purposes (optional)
RUN apk add --no-cache curl

# Expose the port on which the application listens
EXPOSE 8080

# Command to run the executable
CMD ["./picture-api"]
