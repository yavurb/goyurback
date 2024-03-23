package domain

import "context"

type ProjectRepository interface {
	CreateProject(ctx context.Context, project *ProjectCreate) (*Project, error)
}
