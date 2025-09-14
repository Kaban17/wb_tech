# Dockerfile
FROM golang:1.24-alpine

WORKDIR /app

# Устанавливаем зависимости для сборки
RUN apk add --no-cache \
    gcc \
    musl-dev \
    git

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o myapp ./cmd/api

# Экспонируем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./myapp"]
