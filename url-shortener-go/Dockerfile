# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы с зависимостями
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарный файл из builder
COPY --from=builder /app/url-shortener .
# Копируем .env файл если он есть
COPY --from=builder /app/.env .env

EXPOSE 8080

CMD ["./url-shortener"]