package application

import (
	"context"
	"log"

	"github.com/yavurb/goyurback/internal/posts/domain"
)

func (uc *postUsecase) Get(ctx context.Context, id string) (*domain.Post, error) {
	post, err := uc.repository.GetPost(ctx, id)

	if err != nil {
		log.Printf("Error getting post, got error: %v\n", err)
		return nil, domain.ErrPostNotFound
	}

	return post, nil
}
