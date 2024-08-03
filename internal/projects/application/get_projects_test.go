package application

import (
	"context"
	"reflect"
	"testing"

	"github.com/yavurb/goyurback/internal/projects/application/mocks"
	"github.com/yavurb/goyurback/internal/projects/domain"
)

func TestGetProjects(t *testing.T) {
	want := []*domain.Project{
		{
			ID:           1,
			PublicID:     "someid",
			Name:         "Some Project",
			Description:  "Some Description",
			Tags:         []string{"tag1", "tag2"},
			ThumbnailURL: "https://someurl.com/image.jpg",
			WebsiteURL:   "https://somewebsite.com",
			Live:         true,
			PostID:       1,
		},
	}

	repo := &mocks.MockProjectsRepository{
		GetProjectsFn: func(ctx context.Context) ([]*domain.Project, error) { return want, nil },
	}

	uc := NewProjectUsecase(repo)

	projects, err := uc.GetProjects(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(projects, want) {
		t.Errorf("Expected projects to be %v, got: %v", want, projects)
	}
}
