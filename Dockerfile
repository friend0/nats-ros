# Step 1: Use the official Golang image as the builder
FROM golang:1.18 as builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy the Go source code into the container
COPY . .

# Step 4: Download dependencies and build the application
RUN go mod tidy
RUN go build -o nats-ros .

# Step 5: Use a minimal image to run the application
FROM alpine:latest

# Step 6: Set the working directory inside the container
WORKDIR /app

# Step 7: Copy the built binary from the builder container
COPY --from=builder /app/nats-ros .

EXPOSE 4222
EXPOSE 8222
EXPOSE 8080

# Step 9: Specify the command to run the application
CMD ["./nats-ros"]
