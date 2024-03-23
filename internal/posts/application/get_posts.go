package application

import (
	"context"
	"log"

	"github.com/yavurb/goyurback/internal/posts/domain"
)

func (uc *postUsecase) GetPosts(ctx context.Context) ([]*domain.Post, error) {
	posts, err := uc.repository.GetPosts(ctx)

	if err != nil {
		log.Printf("Unable to retrieve posts. Got: %v", err)

		return nil, err
	}

	return posts, nil
}
