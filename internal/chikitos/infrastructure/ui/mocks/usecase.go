package mocks

import (
	"context"

	"github.com/yavurb/goyurback/internal/chikitos/domain"
)

type MockChikitosUsecase struct {
	CreateFn func(ctx context.Context, url, description string) (*domain.Chikito, error)
	GetFn    func(ctx context.Context, id string) (*domain.Chikito, error)
}

func (m *MockChikitosUsecase) Create(ctx context.Context, url, description string) (*domain.Chikito, error) {
	return m.CreateFn(ctx, url, description)
}

func (m *MockChikitosUsecase) Get(ctx context.Context, id string) (*domain.Chikito, error) {
	return m.GetFn(ctx, id)
}
