.PHONY: all setup deps env start db

# Alvo principal
all: setup

# Etapas agrupadas
setup: db deps

# Sobe banco de dados com Docker Compose
db:
	docker-compose up -d

# Instala dependências do Go
deps:
	go mod tidy

# Inicia o aplicativo
start:
	go run .
