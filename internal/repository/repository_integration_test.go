package repository_test

import (
	"context"
	"rate-limiter/internal/repository"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/stretchr/testify/suite"
)

type RedisRepositoryIntegrationTestSuite struct {
	suite.Suite
	client *redis.Client
	repo   repository.RateLimiterRepository
	ctx    context.Context
}

func (s *RedisRepositoryIntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()

	// Conecta ao Redis real em execução (conforme alinhado para usar 127.0.0.1)
	s.client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // sem senha por padrão localmente
		DB:       1,  // Usamos um DB de testes isolado (DB 1) para não sujar o DB principal
	})

	// Valida se o Redis está acessível
	err := s.client.Ping(s.ctx).Err()
	s.Require().NoError(err, "Redis must be running for integration tests")

	s.repo = repository.NewRedisRepository(s.client)
}

func (s *RedisRepositoryIntegrationTestSuite) TearDownTest() {
	// Limpa o DB de testes após cada execução
	s.client.FlushDB(s.ctx)
}

func (s *RedisRepositoryIntegrationTestSuite) TearDownSuite() {
	if s.client != nil {
		s.client.Close()
	}
}

func (s *RedisRepositoryIntegrationTestSuite) TestIncrement_Success() {
	key := "test-ip:127.0.0.1"
	duration := 2 * time.Second

	// Incrementa pela primeira vez
	count, err := s.repo.Increment(s.ctx, key, duration)
	s.NoError(err)
	s.Equal(int64(1), count)

	// Incrementa pela segunda vez
	count, err = s.repo.Increment(s.ctx, key, duration)
	s.NoError(err)
	s.Equal(int64(2), count)
}

func (s *RedisRepositoryIntegrationTestSuite) TestBlockAndIsBlocked() {
	key := "192.168.0.100"
	duration := 2 * time.Second

	// Inicialmente não deve estar bloqueado
	blocked, err := s.repo.IsBlocked(s.ctx, key)
	s.NoError(err)
	s.False(blocked)

	// Aplica o bloqueio
	err = s.repo.Block(s.ctx, key, duration)
	s.NoError(err)

	// Agora deve estar bloqueado
	blocked, err = s.repo.IsBlocked(s.ctx, key)
	s.NoError(err)
	s.True(blocked)
}

func TestRedisRepositoryIntegrationSuite(t *testing.T) {
	suite.Run(t, new(RedisRepositoryIntegrationTestSuite))
}
