# Step 1: Build the binary
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cmd/main .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main cmd/main.go

# Step 2: Use a minimal base image to run the application
FROM alpine:3.19 

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
