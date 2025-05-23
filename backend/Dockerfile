FROM golang:1.24 AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN go build -o server ./cmd/server

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/server .
COPY backend/.env .env

EXPOSE 8080

CMD ["./server"]
