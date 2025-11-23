# Этап сборки (Builder)
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы описания модулей
# (Пути относительны контекста сборки, который будет корнем проекта)
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код проекта
COPY . .

# Собираем бинарник. 
# -o main указывает имя выходного файла
# ./cmd/api указывает, где лежит main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Финальный этап (Runner)
FROM alpine:3.20

WORKDIR /app

# Устанавливаем сертификаты (на всякий случай для внешних запросов)
RUN apk --no-cache add ca-certificates

# Копируем бинарник из этапа сборки
COPY --from=builder /app/main .

# Копируем миграции, чтобы при запуске приложение могло их накатить (если реализуем авто-мигратор)
COPY --from=builder /app/migrations ./migrations

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]