# syntax=docker/dockerfile:1

FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o transaction-service main.go

# 2. Минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/transaction-service .

COPY .env .env


EXPOSE 50052

CMD ["./transaction-service"]