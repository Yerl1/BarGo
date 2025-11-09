FROM golang:1.24-alpine

WORKDIR /app

COPY . .

# ✅ Force build for the same architecture as the container
RUN cd backend && GOOS=linux GOARCH=amd64 go build -o /app/business-web-service ./cmd/business-web-service/main.go

EXPOSE 3001

# ✅ Ensure binary is executable (just in case)
RUN chmod +x /app/business-web-service

CMD ["./business-web-service"]
