# Определение команды Go и флагов
GO = go run
all: server

server:
	sh scripts/server_rebuild.sh

client:
	$(GO) cmd/client/main.go