package application

import (
	"context"
	"errors"
	"log"

	"github.com/yavurb/goyurback/internal/posts/domain"
)

func (uc *postUsecase) Create(ctx context.Context, title, author, slug, description, content string) (*domain.Post, error) {
	postToCreate := &domain.PostCreate{
		Title:       title,
		Author:      author,
		Slug:        slug,
		Description: description,
		Content:     content,
	}

	postCreated, err := uc.repository.CreatePost(ctx, postToCreate)
	if err != nil {
		log.Printf("Error creating post, got error: %v\n", err)
		return nil, errors.New("unable to create post")
	}

	return postCreated, nil
}
