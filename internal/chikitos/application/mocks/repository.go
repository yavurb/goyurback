package mocks

import (
	"context"

	"github.com/yavurb/goyurback/internal/chikitos/domain"
)

type MockChikitosRepository struct {
	CreateChikitoFn func(ctx context.Context, chikito *domain.ChikitoCreate) (*domain.Chikito, error)
	GetChikitoFn    func(ctx context.Context, id string) (*domain.Chikito, error)
}

func (m *MockChikitosRepository) CreateChikito(ctx context.Context, chikito *domain.ChikitoCreate) (*domain.Chikito, error) {
	return m.CreateChikitoFn(ctx, chikito)
}

func (m *MockChikitosRepository) GetChikito(ctx context.Context, id string) (*domain.Chikito, error) {
	return m.GetChikitoFn(ctx, id)
}
