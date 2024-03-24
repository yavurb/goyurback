package application

import (
	"context"
	"log"

	"github.com/yavurb/goyurback/internal/posts/domain"
)

func (uc *postUsecase) Update(ctx context.Context, id string, title, author, slug, description, content string, status domain.Status) (*domain.Post, error) {
	post, err := uc.repository.GetPost(ctx, id)

	if err != nil {
		log.Printf("Error getting post. Got: %v\n", err)

		return nil, domain.ErrPostNotFound
	}

	// TODO: Look for a better solution to 0 values (at the query level maybe)
	if title != "" {
		post.Title = title
	}
	if author != "" {
		post.Author = author
	}
	if slug != "" {
		post.Slug = slug
	}
	if description != "" {
		post.Description = description
	}
	if content != "" {
		post.Content = content
	}
	if status != "" {
		post.Status = status
	}

	postUpdated, err := uc.repository.UpdatePost(ctx, post)

	if err != nil {
		log.Printf("Error updating post. Got: %v\n", err)

		return nil, err
	}

	return postUpdated, nil
}
