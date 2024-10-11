# Step 1: Base image for building Go
FROM golang:1.23-alpine 

WORKDIR /root/backend

# Install air for hot-reloading
RUN go install github.com/air-verse/air@latest

# Copy only the Go module files first
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the backend source code into the container
COPY backend/ .

WORKDIR /root/frontend

# Copy only the Go module files first
COPY frontend/ .

WORKDIR /root/backend

# Expose port 8080
EXPOSE 8080

# Set air as the command to run, it will watch the Go files and reload on changes
CMD ["air"]


