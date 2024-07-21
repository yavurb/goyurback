package mocks

import (
	"context"
	"github.com/yavurb/goyurback/internal/auth/domain"
)

type MockAuthRepository struct {
	CreateAPIKeyFn     func(ctx context.Context, apiKey *domain.APIKeyCreate) (*domain.APIKey, error)
	GetAPIKeyByValueFn func(ctx context.Context, apiKey string) (*domain.APIKey, error)
	RevokeAPIKeyFn     func(ctx context.Context, publicID string) error
}

func (m *MockAuthRepository) CreateAPIKey(ctx context.Context, apiKey *domain.APIKeyCreate) (*domain.APIKey, error) {
	return m.CreateAPIKeyFn(ctx, apiKey)
}
func (m *MockAuthRepository) GetAPIKeyByValue(ctx context.Context, apiKey string) (*domain.APIKey, error) {
	return m.GetAPIKeyByValueFn(ctx, apiKey)
}
func (m *MockAuthRepository) RevokeAPIKey(ctx context.Context, publicID string) error {
	return m.RevokeAPIKeyFn(ctx, publicID)
}
