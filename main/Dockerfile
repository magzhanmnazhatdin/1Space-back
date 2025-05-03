# ─── Builder stage ─────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Устанавливаем git для go modules
RUN apk add --no-cache git

# Копируем файлы с зависимостями и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект и собираем бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# ─── Final stage ──────────────────────────────────────────────────────────────
FROM alpine:latest
WORKDIR /root/

# Копируем бинарник из билдера
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8080

# Указываем точку входа
ENTRYPOINT ["./main"]
