# --- Build stage ---
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# --- Runtime stage ---
FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/main ./main
ENV POLL_INTERVAL_SECONDS=30

CMD ["./main"]