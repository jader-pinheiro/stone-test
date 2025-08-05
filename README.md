# Credito Acordo

`stone-test` Contempla os requisitos do teste realizado para a Stone .

---

## Requisitos

- Go 1.24.4+
---

## Instalação e configuração do Golang
### 1. Instalação do Golang
- Na home baixe o binário do Golang, verifique a versão mais nova no site: https://go.dev/doc/install

   ```bash
  $ sudo rm -rf /usr/local/go && wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz && sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
  ```
- Após ter baixado o binário export o path do Go
     ```bash
    Entre na pasta home $ cd ~ 
    e edite o aquivo .profile
    adicionando a linha no  arquivo:: export PATH=$PATH:/usr/local/go/bin

    atualize com source:
    $ source .profile

    verifique a instalação:
    $ go version

  ```

## Instalação e Execução

### 1. Clone o repositório: https://github.com/jader-pinheiro/stone-test.git
   ```bash
    Dentro do repositório poderá usar os comandos do makefile:

    make → executa docker-compose up -d, go mod tidy, e carrega .env

    make start → roda a aplicação
  ```