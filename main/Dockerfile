#FROM golang:1.23-alpine AS builder
#LABEL authors="Magzhan"
#WORKDIR /app
#
#ENV CGO_ENABLED=0
#ENV GOOS=linux
#ENV GOARCH=amd64
#
## Копирование файлов
#COPY go.mod go.sum ./
#COPY vendor/ ./vendor/
#COPY . .
#
## Сборка приложения
#RUN go build -mod=vendor -ldflags="-s -w" -o ./bin/app ./cmd/app/*
#
## Финальный этап
#FROM scratch
#
## Копируем сертификаты, пользователей и собранный бинарник
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=builder /etc/passwd /etc/passwd
#COPY --from=builder /etc/group /etc/group
#COPY --from=builder /app/bin/app /app/
#
## Установка рабочей директории
#WORKDIR /app
#
## Переключение на непривилегированного пользователя
#USER nobody
#
## Открытие порта
#EXPOSE 8080
#
## Команда запуска приложения
#CMD ["/app/main"]
#
#ENTRYPOINT ["top", "-b"]

# Stage 1: Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Stage 2: Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main /app/main

# Copy SSL certificates (needed for HTTPS requests)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]