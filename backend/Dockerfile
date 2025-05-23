# Stage 1: сборка Go-бинарника (статически, без cgo)
FROM golang:1.24 AS builder
WORKDIR /app

# только манифесты для кеширования
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# исходники
COPY backend/ ./

# отключаем cgo, собираем под Linux/amd64
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -o server ./cmd/server

# Stage 2: минимальный рантайм
FROM debian:bullseye-slim
WORKDIR /app

# копируем собранный бинарь и переменные окружения
COPY --from=builder /app/server .
COPY backend/.env .env

EXPOSE 8080
CMD ["./server"]
