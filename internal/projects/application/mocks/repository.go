package mocks

import (
	"context"

	"github.com/yavurb/goyurback/internal/projects/domain"
)

type MockProjectsRepository struct {
	CreateProjectFn func(ctx context.Context, project *domain.ProjectCreate) (*domain.Project, error)
	GetProjectFn    func(ctx context.Context, id string) (*domain.Project, error)
	GetProjectsFn   func(ctx context.Context) ([]*domain.Project, error)
}

func (m *MockProjectsRepository) CreateProject(ctx context.Context, project *domain.ProjectCreate) (*domain.Project, error) {
	return m.CreateProjectFn(ctx, project)
}

func (m *MockProjectsRepository) GetProject(ctx context.Context, projectID string) (*domain.Project, error) {
	return m.GetProjectFn(ctx, projectID)
}

func (m *MockProjectsRepository) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	return m.GetProjectsFn(ctx)
}
