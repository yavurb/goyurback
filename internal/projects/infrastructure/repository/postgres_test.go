package repository

import (
	"context"
	"errors"
	"log"
	"slices"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
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
	t.Cleanup(func() { connpool.Close() })

	repo := NewRepo(connpool)

	want := &domain.Project{
		ID:           1,
		PublicID:     "pr_18892",
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

	t.Run("CreateProject", func(t *testing.T) {
		t.Run("it should create a project", func(t *testing.T) {
			testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)

			project, err := repo.CreateProject(ctx, &domain.ProjectCreate{
				PublicID:     "pr_18892",
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
			if !project.Compare(*want) {
				t.Errorf("CreateProject() got = %v, want %v", project, want)
			}
		})

		t.Run("it should return an error if we use an unexisting postID", func(t *testing.T) {
			testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)

			project, err := repo.CreateProject(ctx, &domain.ProjectCreate{
				PublicID:     "pr_18892",
				Name:         "Some Project",
				Description:  "Some Description",
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "https://someurl.com/image.jpg",
				WebsiteURL:   "https://somewebsite.com",
				Live:         true,
				PostID:       1,
			})
			if project != nil {
				t.Errorf("CreateProject() got = %v, want nil", project)
			}
			if err == nil {
				t.Error("CreateProject() error = nil, want an error")
			}
			var pgErr *pgconn.PgError
			if !errors.As(err, &pgErr) {
				t.Errorf("CreateProject() error = %v, want a pgconn.PgError", err)
			}
		})
	})

	t.Run("GetProject", func(t *testing.T) {
		t.Run("it should return a project", func(t *testing.T) {
			testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)

			project_, err := repo.CreateProject(ctx, &domain.ProjectCreate{
				PublicID:     "pr_18892",
				Name:         "Some Project",
				Description:  "Some Description",
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "https://someurl.com/image.jpg",
				WebsiteURL:   "https://somewebsite.com",
				Live:         true,
			})
			if err != nil {
				t.Fatal(err)
			}

			project, err := repo.GetProject(ctx, project_.PublicID)
			if err != nil {
				t.Errorf("GetProject() error = %v, want no error", err)
			}
			if project == nil {
				t.Errorf("GetProject() got = %v, want not nil", project)
			}
			if !project.Compare(*want) {
				t.Errorf("GetProject() got = %v, want %v", project, want)
			}
		})
	})

	t.Run("GetProjects", func(t *testing.T) {
		t.Run("it should return a list of projects", func(t *testing.T) {
			testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)
			want := []*domain.Project{
				{
					ID:           1,
					PublicID:     "pr_18892",
					Name:         "Some Project",
					Description:  "Some Description",
					Tags:         []string{"tag1", "tag2"},
					ThumbnailURL: "https://someurl.com/image.jpg",
					WebsiteURL:   "https://somewebsite.com",
					Live:         true,
					PostID:       0,
					CreatedAt:    time.Now().UTC(),
					UpdatedAt:    time.Now().UTC(),
				},
				{
					ID:           2,
					PublicID:     "pr_18893",
					Name:         "Some Project 2",
					Description:  "Some Description 2",
					Tags:         []string{"tag3", "tag4"},
					ThumbnailURL: "https://someurl.com/image2.jpg",
					WebsiteURL:   "https://somewebsite2.com",
					Live:         false,
					PostID:       0,
					CreatedAt:    time.Now().UTC(),
					UpdatedAt:    time.Now().UTC(),
				},
			}
			for _, project := range want {
				_, _ = repo.CreateProject(ctx, &domain.ProjectCreate{
					PublicID:     project.PublicID,
					Name:         project.Name,
					Description:  project.Description,
					Tags:         project.Tags,
					ThumbnailURL: project.ThumbnailURL,
					WebsiteURL:   project.WebsiteURL,
					Live:         project.Live,
				})
			}

			projects, err := repo.GetProjects(ctx)
			if err != nil {
				t.Errorf("GetProjects() error = %v, want no error", err)
			}

			if len(projects) != 2 {
				t.Errorf("GetProjects() got = %d, want 2", len(projects))
			}
			slices.Reverse(want)
			for i, project := range projects {
				if !project.Compare(*want[i]) {
					t.Errorf("GetProjects() got = %v, want %v", project, want[i])
				}
			}
		})
	})
}
