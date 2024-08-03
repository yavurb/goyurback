package mocks

import (
	"context"

	"github.com/yavurb/goyurback/internal/posts/domain"
)

type MockPostsRepository struct {
	CreatePostFn func(ctx context.Context, post *domain.PostCreate) (*domain.Post, error)
	GetPostFn    func(ctx context.Context, id string) (*domain.Post, error)
	GetPostsFn   func(ctx context.Context) ([]*domain.Post, error)
	UpdatePostFn func(ctx context.Context, post *domain.Post) (*domain.Post, error)
}

func (m *MockPostsRepository) CreatePost(ctx context.Context, post *domain.PostCreate) (*domain.Post, error) {
	return m.CreatePostFn(ctx, post)
}

func (m *MockPostsRepository) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	return m.UpdatePostFn(ctx, post)
}

func (m *MockPostsRepository) GetPost(ctx context.Context, postID string) (*domain.Post, error) {
	return m.GetPostFn(ctx, postID)
}

func (m *MockPostsRepository) GetPosts(ctx context.Context) ([]*domain.Post, error) {
	return m.GetPostsFn(ctx)
}
