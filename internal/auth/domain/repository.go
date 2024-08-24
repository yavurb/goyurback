package domain

import "context"

type APIKeyRepository interface {
	CreateAPIKey(ctx context.Context, apiKey *APIKeyCreate) (*APIKey, error)
	GetAPIKeyByValue(ctx context.Context, apiKey string) (*APIKey, error)
	RevokeAPIKey(ctx context.Context, publicID string) error
}
