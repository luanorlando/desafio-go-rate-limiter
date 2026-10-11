# Objetivo
Desenvolver um Rate Limiter em Go que funcione como um middleware para controlar o fluxo de requisições de um serviço web.
O sistema deve ser capaz de limitar o tráfego com base no IP do solicitante ou em um Token de acesso, utilizando o Redis para persistência e orquestração.

# Setup

Para configurar o projeto você pode criar um arquivo .env, já existe um arquivo .env.example de referência.

Configurações disponíveis
|Variável|Descrição|Valor Padrão|
|--------|---------|------------|
|PORT|Porta em que o servidor HTTP irá rodar|8080|
|REDIS_ADDR|Endereço do servidor Redis (host:porta)|localhost:6379 (ou redis:6379 no Docker)|
|PORT|Porta em que o servidor HTTP irá rodar|8080|
|REDIS_DB|Número do banco de dados do Redis|0|
|DEFAULT_LIMIT_IP|Limite máximo de requisições por segundo por IP|10|
|DEFAULT_LIMIT_TOKEN|Limite máximo padrão para Tokens específicos|100|
BLOCK_DURATION_MINUTES|Tempo de bloqueio temporário (em minutos) após estourar o limite|5|

Para validação das regras os valores a serem alterados são
|Variável|Valor Padrão|
|--------|------------|
|DEFAULT_LIMIT_IP|10|
|DEFAULT_LIMIT_TOKEN|100|
|BLOCK_DURATION_MINUTES|5|

> Altere os valores em .env 

Exemplo do trecho com Docker
```go
DEFAULT_LIMIT_IP=4
DEFAULT_LIMIT_TOKEN=10
BLOCK_DURATION_MINUTES=5
```

# Como Executar o Projeto com Docker
Para subir a aplicação e o Redis integrados com um único comando, execute:

```shell
docker compose up --build -d
```

Para derrubar os containers e limpar os recursos:
```shell
docker compose down
```

Como Executar os Testes Automatizados suba o redis
```shell
docker compose up -d redis
# docker compose up --build app-test
```

e depois rode os testes
```shell
go test -v -tags=integration ./...
```

# Como Testar a API
No projeto tem uma pasta `api/` e um arquivo `.http` para testes. Você pode escolher qual usar. As configurações já setadas no projeto é possível testar o caso de sucesso usando qualquer uma das opções, ou pode usar o próprio navegador, nesse exemplo será tudo via terminal

### Sucesso limite por IP
```shell
curl -i http://localhost:8080/
```

> Resultado esperado
```text
HTTP/1.1 200 OK
Date: Sun, 11 Oct 2026 02:59:30 GMT
Content-Length: 41
Content-Type: text/plain; charset=utf-8

Rate Limiter, OK! Requisição permitida.
```

### Sucesso priorizando Token
```shell
curl -i -H "API_KEY: meu-token-super-secreto-123" http://localhost:8080/
```

> Resultado esperado
```text
HTTP/1.1 200 OK
Date: Sun, 11 Oct 2026 03:00:58 GMT
Content-Length: 41
Content-Type: text/plain; charset=utf-8

Rate Limiter, OK! Requisição permitida.
```

## Agora vamos testar validando as regras.
É possível testar segundo fielmente os mesmo valores, porém para uma melhor visualização lembre de diminiur os valores da configurações.

para ficar mais claro, pause sua aplicação e subia a aplicação novamente.
Altere as propriedades para seguir os valores da tabela abaixo do arquivo .env

|Variável|Valor Padrão|
|--------|------------|
|DEFAULT_LIMIT_IP|4|
|DEFAULT_LIMIT_TOKEN|10|
|BLOCK_DURATION_MINUTES|5|

### Sucesso com limite por IP
Executando os request ao mesmo tempo para garantir a validação da regra
```shell
for i in {1..6}; do                
  curl -s -i http://localhost:8080/ &
done
wait
```

> Resultado esperado: 4 resultados com sucesso e 2 com erros

## Sucesso com limite por Token (Token tem preferência)
Nesse exemplo as requisições serão feitas uma após a outra e a regra será validada
```shell
for i in {1..15}; do 
  curl -s -i -H "API_KEY: meu-token-secreto" http://localhost:8080/
  echo ""
done
```

> Resultado esperado 10 sucesso e 5 falhas
Exemplo do erro
```text
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Date: Sun, 11 Oct 2026 03:45:08 GMT
Content-Length: 96

"you have reached the maximum number of requests or actions allowed within a certain time frame"
```

# Regras de Negócio (Critérios de Limitação)
O Rate Limiter deve operar com base em duas estratégias principais:

- [x] **Limitação por IP:** Restringe o número máximo de requisições por segundo recebidas de um único endereço IP.
- [x] **Limitação por Token:** Restringe o número máximo de requisições por segundo baseadas em um token de acesso único.
- [x] **Header:** O token deve ser verificado no header API_KEY: <TOKEN>.
- [x] **Precedência** (Regra de Ouro): As configurações do Token devem se sobrepor às do IP.

> Exemplo: Se o limite global por IP for 10 req/s, mas um token específico tiver permissão para 100 req/s, o sistema deve respeitar o limite do token (100 req/s).
 
# Comportamento em Caso de Bloqueio

Quando o limite de requisições for excedido (seja por IP ou Token):

**Resposta HTTP**: O servidor deve retornar imediatamente o Status Code 429.

**Corpo da Resposta:** Deve retornar exatamente a mensagem: "you have reached the maximum number of requests or actions allowed within a certain time frame"

**Tempo de Bloqueio:** O IP ou Token infrator deve ficar bloqueado por um tempo configurável (ex: 5 minutos). Durante esse período, todas as novas requisições devem ser rejeitadas.
 
# Requisitos Técnicos e Arquitetura

**Middleware:** A lógica do Rate Limiter deve ser implementada como um middleware que envolve o servidor web.

**Persistência (Redis):** As informações de contagem e controle de tempo devem ser armazenadas em um banco de dados Redis (usando imagem Docker).

**Design Pattern (Strategy):**
Você deve implementar o padrão Strategy para a persistência.
Embora o Redis seja obrigatório neste desafio, sua arquitetura deve permitir trocar facilmente o Redis por outro mecanismo de persistência no futuro, apenas mudando a implementação da estratégia.

**Desacoplamento:** A lógica de negócio do limiter deve estar separada da lógica do middleware.

**Configuração:** Todas as definições devem ser feitas via variáveis de ambiente (ou arquivo .env na raiz), incluindo:
Limite máximo de requisições por segundo.
Tempo de bloqueio (em caso de excesso).
Configurações de conexão com Redis.
 
# Entregável
- [x] **Código Fonte:** Repositório contendo a implementação completa.
 **Infraestrutura:**
- [x] **Dockerfile para a aplicação:** docker-compose.yaml que suba a aplicação (na porta 8080) e o banco Redis.
- [x] **Documentação (README):** Explique como configurar o limiter, como alterar as estratégias e como executar o projeto.
- [x] **Testes Automatizados**: Testes que demonstrem a eficácia do limiter e a lógica de precedência (Token > IP).
 
# Regras de Entrega
- [x] **Repositório Único:** O repositório deve conter apenas o projeto em questão.
- [x] **Branch Principal:** Todo o código deve estar na branch main.
- [x] **Execução:** O avaliador deve conseguir rodar o projeto e os testes utilizando apenas o Docker/Docker Compose.
