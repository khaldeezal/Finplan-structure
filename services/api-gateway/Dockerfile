# Dockerfile для Go сервиса
FROM golang:1.24.3  as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/app .

CMD ["./app"]