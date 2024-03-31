package domain

import "context"

type ChikitoRepository interface {
	CreateChikito(ctx context.Context, url, description string) (*Chikito, error)
	GetChikito(ctx context.Context, id string) (*Chikito, error)
}
