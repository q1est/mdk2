# ЭТАП 1: Сборка (Builder)
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл (статический, для работы в минимальном образе)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ЭТАП 2: Финальный образ (Runner)
FROM alpine:latest

# Добавляем сертификаты (обязательно для работы с Neon через SSL)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем только исполняемый файл из первого этапа
COPY --from=builder /app/main .

# Hugging Face требует порт 7860
ENV PORT=7860
EXPOSE 7860

# Запускаем приложение
CMD ["./main"]
