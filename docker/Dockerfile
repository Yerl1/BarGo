# backend/Dockerfile
FROM golang:1.24-alpine

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod ./
RUN go mod download

# Copy backend source
COPY . .

# Build Go binary
RUN go build -o main .

# Expose port
EXPOSE 8080

# Run the server
CMD ["./main"]
