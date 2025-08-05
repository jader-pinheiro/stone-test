.PHONY: all setup deps start db

# Alvo principal
all: setup

# Etapas agrupadas
setup: db deps

# Remove container antigo e volumes (dados), e sobe o banco de dados limpo
db:
	docker-compose down -v
	docker-compose up -d

# Instala dependências do Go
deps:
	go mod tidy

# Remove container antigo, sobe banco e roda app
start:
	go run .
