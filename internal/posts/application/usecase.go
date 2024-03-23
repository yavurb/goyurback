package application

import "github.com/yavurb/goyurback/internal/posts/domain"

type postUsecase struct {
	repository domain.PostRepository
}

func NewPostUsecase(repository domain.PostRepository) domain.PostUsecase {
	return &postUsecase{repository}
}
