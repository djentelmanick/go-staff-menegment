# Билд стадии
FROM golang:1.21 as builder

WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/

# Финальная стадия
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из builder стадии
COPY --from=builder /app/main .
# Копируем статические файлы
COPY static ./static

EXPOSE 8080

CMD ["./main"]