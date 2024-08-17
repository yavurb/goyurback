package repository

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavurb/goyurback/internal/projects/domain"
	"github.com/yavurb/goyurback/testhelpers"
)

func TestPostgresRepository(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(t, ctx)
	if err != nil {
		t.Errorf("Error creating container: %s", err)
	}

	connpool, err := pgxpool.New(ctx, pgContainer.ConnString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	repo := NewRepo(connpool)

	t.Run("CreateProject", func(t *testing.T) {
		t.Run("it should create a project", func(t *testing.T) {
			want := &domain.Project{
				ID:           1,
				PublicID:     "someid",
				Name:         "Some Project",
				Description:  "Some Description",
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "https://someurl.com/image.jpg",
				WebsiteURL:   "https://somewebsite.com",
				Live:         true,
				PostID:       0,
				CreatedAt:    time.Now().UTC(),
				UpdatedAt:    time.Now().UTC(),
			}
			project, err := repo.CreateProject(ctx, &domain.ProjectCreate{
				Name:         "Some Project",
				Description:  "Some Description",
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "https://someurl.com/image.jpg",
				WebsiteURL:   "https://somewebsite.com",
				Live:         true,
			})
			if err != nil {
				t.Errorf("CreateProject() error = %v, want no error", err)
			}
			if project == nil {
				t.Errorf("CreateProject() got = %v, want not nil", project)
			}

			// FIXME: These fields shouldn't be set here
			want.PublicID = project.PublicID
			want.CreatedAt = project.CreatedAt
			want.UpdatedAt = project.UpdatedAt

			if !reflect.DeepEqual(project, want) {
				t.Errorf("CreateProject() got = %v, want %v", project, want)
			}
		})
	})
}
