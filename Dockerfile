# Step 1: Build the Go backend app
FROM golang:1.23-alpine AS builder

# Set the working directory for the backend code inside the container
WORKDIR /app/backend

# Copy only the Go module files and vendor folder (if applicable) first
COPY backend/go.mod backend/go.sum backend/vendor ./

# Download dependencies if necessary
# RUN go mod download

# Copy the rest of the backend source code into the container
COPY backend/ .

# Build the Go binary, using the vendor directory if available
RUN go build -mod=vendor -o go-server .

# Step 2: Create the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled Go backend binary from the builder stage
COPY --from=builder /app/backend/go-server .

# Copy the entire backend folder to maintain the structure (optional if you need to preserve the folder structure)
COPY backend/ /root/backend/

# Copy the frontend folder as-is to maintain the structure
COPY frontend/ /root/frontend/

# Expose port 8080 for the Go backend
EXPOSE 80

# Run the Go server binary
CMD ["./go-server"]
