FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN cd backend && go build -o /app/auth-service ./cmd/auth-service/main.go

EXPOSE 3000

CMD ["./auth-service"]
