FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN go build -o auth-service ./backend/cmd/auth-service/main.go

EXPOSE 3000

CMD ["./auth-service", "--mode=au"]