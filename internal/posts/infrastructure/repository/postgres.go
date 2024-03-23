package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavurb/goyurback/internal/database/postgres"
	"github.com/yavurb/goyurback/internal/pgk/publicid"
	"github.com/yavurb/goyurback/internal/posts/domain"
)

const prefix = "po"

type Repository struct {
	db *postgres.Queries
}

func NewRepo(connpool *pgxpool.Pool) domain.PostRepository {
	db := postgres.New(connpool)

	return &Repository{
		db: db,
	}
}

func (r *Repository) CreatePost(ctx context.Context, post *domain.PostCreate) (*domain.Post, error) {
	id, _ := publicid.New(prefix) // TODO: handle errors and validate if the id already exists

	post_, err := r.db.CreatePost(ctx, postgres.CreatePostParams{
		PublicID:    id,
		Title:       post.Title,
		Author:      post.Author,
		Slug:        post.Slug,
		Description: post.Description,
		Content:     post.Content,
	})

	if err != nil {
		return nil, err
	}

	newPost := &domain.Post{
		ID:          post_.ID,
		PublicID:    post_.PublicID,
		Title:       post_.Title,
		Author:      post_.Author,
		Slug:        post_.Slug,
		Description: post_.Description,
		Content:     post_.Content,
		Status:      domain.Status(post_.Status),
		PublishedAt: post_.PublishedAt.Time,
		CreatedAt:   post_.CreatedAt.Time,
		UpdatedAt:   post_.UpdatedAt.Time,
	}

	return newPost, nil
}
