package domain

import "context"

type PostRepository interface {
	// GetPost(ctx context.Context, id int) (*Post, error)
	// GetPosts(ctx context.Context) ([]*Post, error)
	// GetPostBySlug(ctx context.Context, slug string) (*Post, error)
	CreatePost(ctx context.Context, post *PostCreate) (*Post, error)
	// UpdatePost(ctx context.Context, post *Post) (*Post, error)
}
