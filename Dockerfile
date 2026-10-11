# --- Stage 1: Build da aplicação ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependências primeiro para aproveitar o cache do Docker
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código fonte para dentro do container
COPY . .

# Compila o binário apontando para a pasta cmd/main.go
# CGO_ENABLED=0 garante um binário estático para rodar no Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/rate-limiter ./cmd/main.go

# --- Stage 2: Imagem final leve para execução ---
FROM alpine:latest

WORKDIR /root/

# Copia o binário compilado do estágio anterior
COPY --from=builder /app/rate-limiter .

# Copia o arquivo .env se você o utiliza em produção/desenvolvimento no container
COPY .env .env

# Expõe a porta em que a aplicação roda (ex: 8080)
EXPOSE 8080

# Comando para iniciar a aplicação
CMD ["./rate-limiter"]