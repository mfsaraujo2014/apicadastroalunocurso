# API Cadastro de Alunos e Cursos

Este é o repositório da API para cadastro de alunos e cursos. Aqui você encontrará o código-fonte da API desenvolvida em Go.

## Pré-requisitos

Antes de começar, certifique-se de ter instalado em sua máquina:

- Go 
- Docker
- Docker Compose

## Como Executar Localmente

Siga estas etapas para configurar e executar o projeto em sua máquina local:

1. **Clone este repositório para sua máquina local:**

   ```bash
   git clone https://github.com/seu-usuario/apicadastroalunocurso.git

2. **Crie um arquivo .env na raiz do projeto e configure as variáveis de ambiente necessárias:**

    ```bash
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=sua-senha-secreta
    DB_NAME=postgres
    API_PORT=9001

3. **Execute o seguinte comando para instalar as dependências do projeto:**

    ```bash
    go mod tidy

4. **Inicie o banco de dados PostgreSQL em um Container Docker:**

    ```bash
    docker-compose up -d
    ou
    make run

5. **Execute o projeto Go:**

    ```bash
    go run main.go

**Uso da API**

A API estará disponível em http://localhost:9001. Consulte a documentação da API para obter informações sobre os endpoints disponíveis e como usá-los.

**Contribuindo**

Contribuições são bem-vindas! Sinta-se à vontade para abrir uma issue ou enviar um pull request.
