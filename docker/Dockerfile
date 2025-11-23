FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем ВЕСЬ код (включая internal)
COPY . .

# Собираем
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Финальный образ
FROM alpine:3.20

WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]