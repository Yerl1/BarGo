FROM golang:1.24-alpine

WORKDIR /app

# 1) ставим curl для healthcheck
RUN apk add --no-cache curl

# 2) копируем проект
COPY . .

# 3) собираем бинарь
RUN cd backend && go build -o /app/auth-service ./cmd/auth-service/main.go

EXPOSE 3004

CMD ["./auth-service"]
