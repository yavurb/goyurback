package domain

import (
	"context"
)

type ChikitoUsecase interface {
	Create(ctx context.Context, url, description string) (*Chikito, error)
}
