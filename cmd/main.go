package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rate-limiter/internal/config"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/middleware"
	"rate-limiter/internal/repository"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Erro ao conectar com o Redis: %v", err)
	}
	log.Println("Conectado ao Redis com sucesso!")

	repo := repository.NewRedisRepository(rdb)

	rateLimiter := limiter.NewRateLimiter(
		repo,
		cfg.DefaultLimitIP,
		cfg.DefaultLimitToken,
		cfg.BlockDuration,
	)

	mw := middleware.NewRateLimiterMiddleware(rateLimiter)

	mux := http.NewServeMux()

	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Rate Limiter, OK! Requisição permitida."))
	})

	mux.Handle("/", mw.Handle(helloHandler))

	portStr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Servidor rodando na porta %s...", cfg.Port)
	if err := http.ListenAndServe(portStr, mux); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

/*
Requisição normal (Permitida):

curl -i http://localhost:8080/
Resposta esperada: Status 200 OK com a mensagem "Bem-vindo ao Rate Limiter! Requisição permitida.".

Estourar o limite por IP (Teste do Rate Limiter):
Como o limite padrão de IP está configurado para 10 requisições por segundo, vamos disparar 15 requisições de uma vez só:

for i in {1..15}; do
  curl -s -i http://localhost:8080/
  echo ""
done
Resultado esperado: As primeiras 10 passam (200 OK)
HTTP/1.1 200 OK
Date: Sat, 10 Oct 2026 21:23:18 GMT
Content-Length: 41
Content-Type: text/plain; charset=utf-8

e as seguintes retornam o Status 429 Too Many Requests com a mensagem exigida pelo desafio.
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Date: Sat, 10 Oct 2026 21:25:53 GMT
Content-Length: 96

"you have reached the maximum number of requests or actions allowed within a certain time frame"

Testar com Token (API_KEY):

for i in {1..15}; do
  curl -s -i -H "API_KEY: meu-token-secreto" http://localhost:8080/
  echo ""
done
Resultado esperado: Como o token tem um limite maior (ex: 100 por segundo), ele será permitido (200 OK).
Rate Limiter, OK! Requisição permitida.
HTTP/1.1 200 OK
Date: Sat, 10 Oct 2026 21:28:24 GMT
Content-Length: 41
Content-Type: text/plain; charset=utf-8

para teste negativo altere o .env a variavel DEFAULT_LIMIT_TOKEN, altere para um número menor que 10
Lembre de rodar a aplicação novamente(Se subiu tudo via docker, terá que rodar subir as imagens novamente para pegar o novo valor das variaveis de .env) e o resultado será
*/
