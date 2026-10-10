package main

func main() {

}

/*
meu-rate-limiter/
├── cmd/
│   └── server/
│       └── main.go          # Ponto de entrada da aplicação
├── internal/
│   ├── config/              # Leitura de variáveis de ambiente (.env)
│   ├── limiter/             # Lógica de negócio do Rate Limiter
│   ├── middleware/          # Middleware HTTP que intercepta as requisições
│   └── repository/          # Camada de persistência (com Strategy para o Redis)
├── .env                     # Variáveis de ambiente
├── Dockerfile
├── docker-compose.yaml
└── README.md
*/
