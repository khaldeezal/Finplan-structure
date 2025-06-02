#  Билдим бинарник вне контейнера (multi-stage)
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-service .

COPY .env .env

EXPOSE 50051

CMD ["./auth-service"]