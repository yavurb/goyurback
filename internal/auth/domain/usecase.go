package domain

import "context"

type APIKeyUsecase interface {
	CreateAPIKey(ctx context.Context, name string) (*APIKey, error)
	RevokeAPIKey(ctx context.Context, publicID string) error
}
