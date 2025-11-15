FROM golang:1.24-alpine

WORKDIR /app

COPY . .

# ✅ Force build for the same architecture as the container
RUN cd backend && go build -o /app/consumer-web-service ./cmd/consumer-web-service/main.go

EXPOSE 3000

# ✅ Ensure binary is executable (just in case)
# RUN chmod +x /app/consumer-web-service

CMD ["./consumer-web-service"]
