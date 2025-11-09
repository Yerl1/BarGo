FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN cd backend && go build -o /app/business-service ./cmd/business-service/main.go

EXPOSE 3003

CMD ["./business-service"]
