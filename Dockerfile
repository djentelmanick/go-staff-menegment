# Используем официальный образ Go
FROM golang:1.21-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Используем минимальный образ для продакшена
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Создаем рабочую директорию
WORKDIR /root/

# Копируем исполняемый файл из builder
COPY --from=builder /app/main .

# Копируем статические файлы
COPY --from=builder /app/static ./static

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]