package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

	// TODO: print log
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

func (r *Repository) GetPost(ctx context.Context, id string) (*domain.Post, error) {
	post_, err := r.db.GetPost(ctx, id)

	if err != nil {
		log.Printf("Error getting post. Got: %v", err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}

		log.Panicf("DB Error obtaining post: %v", err)
	}

	post := &domain.Post{
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

	return post, nil
}

// TODO: Add pagination
func (r *Repository) GetPosts(ctx context.Context) ([]*domain.Post, error) {
	posts, err := r.db.GetPosts(ctx)

	if err != nil {
		log.Printf("DB Error obtaining posts: %v", err)

		if errors.Is(err, pgx.ErrNoRows) {
			return []*domain.Post{}, nil
		}

		return nil, err
	}

	posts_ := []*domain.Post{}

	for _, post := range posts {
		posts_ = append(posts_, &domain.Post{
			ID:          post.ID,
			PublicID:    post.PublicID,
			Title:       post.Title,
			Author:      post.Author,
			Slug:        post.Slug,
			Description: post.Description,
			Content:     post.Content,
			Status:      domain.Status(post.Status),
			PublishedAt: post.PublishedAt.Time,
			CreatedAt:   post.CreatedAt.Time,
			UpdatedAt:   post.UpdatedAt.Time,
		})
	}

	return posts_, nil
}

func (r *Repository) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	post_, err := r.db.UpdatePost(ctx, postgres.UpdatePostParams{
		ID:          post.ID,
		Title:       post.Title,
		Author:      post.Author,
		Slug:        post.Slug,
		Description: post.Description,
		Content:     post.Content,
		Status:      postgres.PostStatus(post.Status),
		PublishedAt: pgtype.Timestamp{
			Valid: true,
			Time:  post.PublishedAt,
		},
	})

	if err != nil {
		log.Printf("DB Error updating post: %v\n", err)

		return nil, err
	}

	postUpdated := &domain.Post{
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

	return postUpdated, nil
}
