.PHONY: all setup deps env start db

# Diretório do projeto
REPO_DIR=stone-test

# Alvo principal
all: setup

# Etapas agrupadas
setup: db deps env

# Sobe banco de dados com Docker Compose
db:
	docker-compose up -d

# Instala dependências do Go
deps:
	cd $(REPO_DIR) && go mod tidy

# Carrega variáveis de ambiente do .env
env:
	cd $(REPO_DIR) && set -a && source .env && set +a

# Inicia o aplicativo
start:
	cd $(REPO_DIR) && go run main.go
