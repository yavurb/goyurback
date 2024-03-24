package domain

import "context"

type ProjectRepository interface {
	CreateProject(ctx context.Context, project *ProjectCreate) (*Project, error)
	GetProject(ctx context.Context, id string) (*Project, error)
	GetProjects(ctx context.Context) ([]*Project, error)
}
