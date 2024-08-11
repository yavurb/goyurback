package application

import (
	"context"
	"log"

	"github.com/yavurb/goyurback/internal/posts/domain"
)

func (uc *postUsecase) Update(ctx context.Context, id string, title, author, slug, description, content *string, status *domain.Status) (*domain.Post, error) {
	post, err := uc.repository.GetPost(ctx, id)
	if err != nil {
		log.Printf("Error getting post. Got: %v\n", err)

		return nil, domain.ErrPostNotFound
	}

	if title != nil {
		post.Title = *title
	}
	if author != nil {
		post.Author = *author
	}
	if slug != nil {
		post.Slug = *slug
	}
	if description != nil {
		post.Description = *description
	}
	if content != nil {
		post.Content = *content
	}
	if status != nil {
		post.Status = *status
	}

	postUpdated, err := uc.repository.UpdatePost(ctx, post)
	if err != nil {
		log.Printf("Error updating post. Got: %v\n", err)

		return nil, err
	}

	return postUpdated, nil
}
