FROM golang:1.24-alpine

WORKDIR /app

COPY . .

# ✅ Force build for the same architecture as the container
RUN cd backend && GOOS=linux GOARCH=amd64 go build -o /app/web-service ./cmd/web-service/main.go

EXPOSE 8080

# ✅ Ensure binary is executable (just in case)
RUN chmod +x /app/web-service

CMD ["./web-service"]
