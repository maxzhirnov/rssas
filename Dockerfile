# Этап сборки
FROM golang:1.20 as builder

WORKDIR /app

# Копирование модулей и их установка
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода приложения
COPY . .

# Сборка cmd/client/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/rssapp ./cmd/main.go


# Этап выполнения
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=builder /app/rssapp /rssapp

# Запустите приложение
ENTRYPOINT ["/rssapp"]