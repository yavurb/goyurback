package mocks

import (
	"context"

	"github.com/yavurb/goyurback/internal/posts/domain"
)

type MockPostsUsecase struct {
	GetFn      func(ctx context.Context, id string) (*domain.Post, error)
	GetPostsFn func(ctx context.Context) ([]*domain.Post, error)
	CreateFn   func(ctx context.Context, title, author, slug, description, content string) (*domain.Post, error)
	UpdateFn   func(ctx context.Context, id string, title, author, slug, description, content *string, status *domain.Status) (*domain.Post, error)
}

func (m *MockPostsUsecase) Get(ctx context.Context, id string) (*domain.Post, error) {
	return m.GetFn(ctx, id)
}

func (m *MockPostsUsecase) GetPosts(ctx context.Context) ([]*domain.Post, error) {
	return m.GetPostsFn(ctx)
}

func (m *MockPostsUsecase) Create(ctx context.Context, title, author, slug, description, content string) (*domain.Post, error) {
	return m.CreateFn(ctx, title, author, slug, description, content)
}

func (m *MockPostsUsecase) Update(ctx context.Context, id string, title, author, slug, description, content *string, status *domain.Status) (*domain.Post, error) {
	return m.UpdateFn(ctx, id, title, author, slug, description, content, status)
}
