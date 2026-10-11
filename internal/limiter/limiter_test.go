package limiter_test

import (
	"context"
	"testing"
	"time"

	"rate-limiter/internal/limiter"

	"github.com/stretchr/testify/suite"
)

type MockRepository struct {
	storage map[string]int64
	blocked map[string]time.Time
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		storage: make(map[string]int64),
		blocked: make(map[string]time.Time),
	}
}

func (m *MockRepository) IsBlocked(ctx context.Context, key string) (bool, error) {
	if blockUntil, exists := m.blocked[key]; exists {
		if time.Now().Before(blockUntil) {
			return true, nil
		}
		delete(m.blocked, key)
	}
	return false, nil
}

func (m *MockRepository) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	m.storage[key]++
	return m.storage[key], nil
}

func (m *MockRepository) Block(ctx context.Context, key string, duration time.Duration) error {
	m.blocked[key] = time.Now().Add(duration)
	return nil
}

type RateLimiterTestSuite struct {
	suite.Suite
	repo *MockRepository
	ctx  context.Context
}

func (s *RateLimiterTestSuite) SetupTest() {
	s.repo = NewMockRepository()
	s.ctx = context.Background()
}

func (s *RateLimiterTestSuite) TestAllow_UnderLimit() {
	l := limiter.NewRateLimiter(s.repo, 2, 5, time.Minute)

	allowed, err := l.Allow(s.ctx, "192.168.1.1", "")
	s.NoError(err)
	s.True(allowed, "First request should be allowed")

	allowed, err = l.Allow(s.ctx, "192.168.1.1", "")
	s.NoError(err)
	s.True(allowed, "Second request should be allowed")
}

func (s *RateLimiterTestSuite) TestAllow_ExceededLimit() {
	l := limiter.NewRateLimiter(s.repo, 1, 5, time.Minute)

	allowed, err := l.Allow(s.ctx, "192.168.1.2", "")
	s.NoError(err)
	s.True(allowed)

	allowed, err = l.Allow(s.ctx, "192.168.1.2", "")
	s.NoError(err)
	s.False(allowed, "Request exceeding limit should be blocked")
}

func (s *RateLimiterTestSuite) TestTokenPrecedenceOverIP() {
	l := limiter.NewRateLimiter(s.repo, 1, 5, time.Minute)

	l.Allow(s.ctx, "192.168.1.3", "")
	allowedIP, _ := l.Allow(s.ctx, "192.168.1.3", "")
	s.False(allowedIP, "IP should be blocked on second attempt")

	allowedToken, err := l.Allow(s.ctx, "192.168.1.3", "meu-token-valido")
	s.NoError(err)
	s.True(allowedToken, "Token request must be allowed due to precedence even if IP is blocked")
}

func TestRateLimiterSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}
