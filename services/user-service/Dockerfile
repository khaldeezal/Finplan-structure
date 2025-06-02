# syntax=docker/dockerfile:1

FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service main.go

# 2. Минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-service .

COPY .env .env


EXPOSE 50053

CMD ["./user-service"]