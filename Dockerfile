# Этап, на котором выполняется сборка приложения
FROM golang:1.22-alpine as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /main app/main.go
# Финальный этап, копируем собранное приложение
FROM alpine:3
COPY --from=builder main /bin/main

EXPOSE 8000:8000

ENTRYPOINT ["/bin/main"]