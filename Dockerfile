# ЭТАП 1: Сборка (Builder)
FROM golang:1.24-alpine AS builder
ENV GOTOOLCHAIN=auto
# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы зависимостей из КОРНЯ (так как папки back нет)
COPY go.mod go.sum ./
RUN go mod download

# Копируем все остальные файлы проекта (включая main.go) в рабочую директорию
COPY . .

# Собираем статический бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ЭТАП 2: Финальный образ (Runner)
FROM alpine:latest

# Устанавливаем сертификаты для безопасного подключения к базе данных (Neon)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем скомпилированный файл из этапа сборки
COPY --from=builder /app/main .

# Настройки для Hugging Face Spaces
ENV PORT=7860
EXPOSE 7860

# Запускаем приложение
CMD ["./main"]
