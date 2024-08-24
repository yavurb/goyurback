package domain

import (
	"context"
)

type PostUsecase interface {
	Get(ctx context.Context, id string) (*Post, error)
	GetPosts(ctx context.Context) ([]*Post, error)
	// GetBySlug(slug string) (*Post, error)
	Create(ctx context.Context, title, author, slug, description, content string) (*Post, error)
	Update(ctx context.Context, id string, title, author, slug, description, content *string, status *Status) (*Post, error)
}
