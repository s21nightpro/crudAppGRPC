# Используйте официальный образ Go как базовый
FROM golang:1.17-alpine as builder

# Установите рабочую директорию в контейнере
WORKDIR /app

# Копируйте go.mod и go.sum в контейнер
COPY go.mod go.sum./

# Загрузите зависимости
RUN go mod download

# Копируйте исходный код в контейнер
COPY..

# Соберите приложение
RUN go build -o main.

# Используйте официальный образ Alpine Linux как базовый для финального образа
FROM alpine:latest

# Установите рабочую директорию в контейнере
WORKDIR /root/

# Копируйте исполняемый файл из предыдущего шага
COPY --from=builder /app/main.

# Установите переменную окружения для запуска приложения
ENTRYPOINT ["./main"]
