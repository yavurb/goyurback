package domain

import "context"

type ProjectUsecase interface {
	Create(ctx context.Context, name, description, thumbnailURL, websiteURL string, live bool, tags []string, postId int32) (*Project, error)
	Get(ctx context.Context, id string) (*Project, error)
	GetProjects(ctx context.Context) ([]*Project, error)
}
