package domain

import "context"

type PostUsecase interface {
	Get(ctx context.Context, id string) (*Post, error)
	// GetPosts() ([]*Post, error)
	// GetBySlug(slug string) (*Post, error)
	Create(ctx context.Context, title, author, slug, description, content string) (*Post, error)
	// Update(id int, post *PostUpdate) (*Post, error)
}
