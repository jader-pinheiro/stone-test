# Stone Teste

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

Após clonar o repositório baixe os arquivos de stokes descompacte ele, os arquivos baixados estarão em .txt e deverão permanecer nessa extensão pois o processamento está para arquivos .txt e coloque eles dentro da pasta `file` que se encontra na raiz do projeto
com o comando make start a aplicação irá processar os arquivos fazendo a inserção no banco de dados e após o término
a API de consulta estará disponível no seguinte endpoint:

``localhost:3000/ticker/DIIF30F31?startDate=2025-07-28``

Conforme solicitado pela documentação do teste o parametro startDate não é obrigatório
   ```bash
    Dentro do repositório execute os comandos do makefile nessa ordem:

    make db → sobe o container do Postgres
    make deps → baixa as dependencias do projeto
    make start → roda a aplicação em golang, 
  ```