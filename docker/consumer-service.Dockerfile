FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN cd backend && go build -o /app/consumer-service ./cmd/consumer-service/main.go

EXPOSE 3002

CMD ["./consumer-service"]
