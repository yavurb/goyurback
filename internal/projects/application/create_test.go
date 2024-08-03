package application

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/yavurb/goyurback/internal/projects/application/mocks"
	"github.com/yavurb/goyurback/internal/projects/domain"
)

func TestCreateProject(t *testing.T) {
	want := &domain.Project{
		ID:           1,
		PublicID:     "someid",
		Name:         "Some Project",
		Description:  "Some Description",
		Tags:         []string{"tag1", "tag2"},
		ThumbnailURL: "https://someurl.com/image.jpg",
		WebsiteURL:   "https://somewebsite.com",
		Live:         true,
		PostID:       1,
	}

	repo := &mocks.MockProjectsRepository{
		CreateProjectFn: func(ctx context.Context, project *domain.ProjectCreate) (*domain.Project, error) {
			return &domain.Project{
				ID:           1,
				PublicID:     "someid",
				Name:         project.Name,
				Description:  project.Description,
				Tags:         project.Tags,
				ThumbnailURL: project.ThumbnailURL,
				WebsiteURL:   project.WebsiteURL,
				Live:         project.Live,
				PostID:       project.PostID,
			}, nil
		},
	}

	uc := NewProjectUsecase(repo)
	ctx := context.Background()

	project, err := uc.Create(ctx, want.Name, want.Description, want.ThumbnailURL, want.WebsiteURL, want.Live, want.Tags, want.PostID)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(project, want) {
		t.Errorf("Expected project to be %v, got: %v", want, project)
	}
}

func TestCreateProjectWithDBError(t *testing.T) {
	want := "unable to create project"
	repo := &mocks.MockProjectsRepository{
		CreateProjectFn: func(ctx context.Context, project *domain.ProjectCreate) (*domain.Project, error) {
			return nil, errors.New("DB Error")
		},
	}

	uc := NewProjectUsecase(repo)
	ctx := context.Background()

	_, err := uc.Create(ctx, "Some Project", "Some Description", "https://someurl.com/image.jpg", "https://somewebsite.com", true, []string{"tag1", "tag2"}, 1)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if err.Error() != want {
		t.Errorf("Expected error to be %v, got: %v", want, err)
	}
}
