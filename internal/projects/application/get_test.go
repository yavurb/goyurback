package application

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/yavurb/goyurback/internal/projects/application/mocks"
	"github.com/yavurb/goyurback/internal/projects/domain"
)

func TestGetProject(t *testing.T) {
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
		GetProjectFn: func(ctx context.Context, id string) (*domain.Project, error) {
			return &domain.Project{
				ID:           1,
				PublicID:     id,
				Name:         want.Name,
				Description:  want.Description,
				Tags:         want.Tags,
				ThumbnailURL: want.ThumbnailURL,
				WebsiteURL:   want.WebsiteURL,
				Live:         want.Live,
				PostID:       want.PostID,
			}, nil
		},
	}

	uc := NewProjectUsecase(repo)

	project, err := uc.Get(context.Background(), want.PublicID)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if project == nil {
		t.Errorf("Expected project to be %v, got nil", want)
	}

	if !reflect.DeepEqual(project, want) {
		t.Errorf("Expected project to be %v, got: %v", want, project)
	}
}

func TestGetProjectNotFound(t *testing.T) {
	repo := &mocks.MockProjectsRepository{
		GetProjectFn: func(ctx context.Context, id string) (*domain.Project, error) {
			return nil, errors.New("project not found")
		},
	}

	uc := NewProjectUsecase(repo)
	project, err := uc.Get(context.Background(), "someid")

	if !errors.Is(err, domain.ErrProjectNotFound) {
		t.Errorf("Expected error to be %v, got: %v", domain.ErrProjectNotFound, err)
	}

	if project != nil {
		t.Errorf("Expected project to be nil, got: %v", project)
	}
}
