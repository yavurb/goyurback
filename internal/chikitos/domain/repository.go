package domain

import "context"

type ChikitoRepository interface {
	CreateChikito(ctx context.Context, chikito *ChikitoCreate) (*Chikito, error)
	GetChikito(ctx context.Context, id string) (*Chikito, error)
}
