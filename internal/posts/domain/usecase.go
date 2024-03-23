package domain

import "context"

type PostUsecase interface {
	// Get(id int) (*Post, error)
	// GetPosts() ([]*Post, error)
	// GetBySlug(slug string) (*Post, error)
	Create(ctx context.Context, title, author, slug, description, content string) (*Post, error)
	// Update(id int, post *PostUpdate) (*Post, error)
}
