# Определение команды Go и флагов
GO = go
GOFLAGS = -o
TARGET = crudApp

# Целевая задача по умолчанию: сборка исполняемого файла
all: server

# Правила для сборки целевого исполняемого файла
server:
	sh scripts/server_rebuild.sh

client:
	go run cmd/client/main.go