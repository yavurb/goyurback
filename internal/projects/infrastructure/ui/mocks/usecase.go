package mocks

import (
	"context"

	"github.com/yavurb/goyurback/internal/projects/domain"
)

type MockProjectsUsecase struct {
	CreateFn      func(ctx context.Context, name, description, thumbnailURL, websiteURL string, live bool, tags []string, postId int32) (*domain.Project, error)
	GetFn         func(ctx context.Context, id string) (*domain.Project, error)
	GetProjectsFn func(ctx context.Context) ([]*domain.Project, error)
}

func (uc *MockProjectsUsecase) Create(ctx context.Context, name, description, thumbnailURL, websiteURL string, live bool, tags []string, postId int32) (*domain.Project, error) {
	return uc.CreateFn(ctx, name, description, thumbnailURL, websiteURL, live, tags, postId)
}

func (uc *MockProjectsUsecase) Get(ctx context.Context, id string) (*domain.Project, error) {
	return uc.GetFn(ctx, id)
}

func (uc *MockProjectsUsecase) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	return uc.GetProjectsFn(ctx)
}
