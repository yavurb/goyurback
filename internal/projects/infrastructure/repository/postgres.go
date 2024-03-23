package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavurb/goyurback/internal/database/postgres"
	"github.com/yavurb/goyurback/internal/pgk/publicid"
	"github.com/yavurb/goyurback/internal/projects/domain"
)

const prefix = "pr"

type Repository struct {
	db *postgres.Queries
}

func NewRepo(connpool *pgxpool.Pool) domain.ProjectRepository {
	return &Repository{
		db: postgres.New(connpool),
	}
}

func (r *Repository) CreateProject(ctx context.Context, project *domain.ProjectCreate) (*domain.Project, error) {
	var postID pgtype.Int4 = pgtype.Int4{Valid: false}

	if _postID := project.PostID; _postID != 0 {
		postID = pgtype.Int4{Int32: project.PostID, Valid: true}
	}

	publicID, _ := publicid.New(prefix) // TODO: Handle error and add a retry mechanism to validate if the id already exists

	project_, err := r.db.CreateProject(ctx, postgres.CreateProjectParams{
		PublicID:     publicID,
		Name:         project.Name,
		Description:  project.Description,
		Tags:         project.Tags,
		ThumbnailUrl: project.ThumbnailURL,
		WebsiteUrl:   project.WebsiteURL,
		Live:         project.Live,
		PostID:       postID,
	})

	if err != nil {
		return nil, err
	}

	newProject := &domain.Project{
		ID:           project_.ID,
		PublicID:     project_.PublicID,
		Name:         project_.Name,
		Description:  project_.Description,
		Tags:         project_.Tags,
		ThumbnailURL: project_.ThumbnailUrl,
		WebsiteURL:   project_.WebsiteUrl,
		Live:         project_.Live,
		PostID:       project_.PostID.Int32,
		CreatedAt:    project_.CreatedAt.Time,
		UpdatedAt:    project_.UpdatedAt.Time,
	}

	return newProject, nil
}
