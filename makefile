.PHONY: all setup deps env start db

# Alvo principal
all: setup

# Etapas agrupadas
setup: db deps env

# Sobe banco de dados com Docker Compose
db:
	docker-compose up -d

# Instala depend�ncias do Go
deps:
	go mod tidy

# Carrega vari�veis de ambiente do .env
env:
	set -a && source .env && set +a

# Inicia o aplicativo
start:
	go run main.go
