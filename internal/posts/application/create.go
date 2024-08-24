package application

import (
	"context"
	"errors"
	"log"

	"github.com/yavurb/goyurback/internal/pgk/ids"
	"github.com/yavurb/goyurback/internal/posts/domain"
)

const prefix = "po"

func (uc *postUsecase) Create(ctx context.Context, title, author, slug, description, content string) (*domain.Post, error) {
	id, _ := ids.NewPublicID(prefix) // TODO: handle errors and validate if the id already exists

	postToCreate := &domain.PostCreate{
		PublicID:    id,
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
