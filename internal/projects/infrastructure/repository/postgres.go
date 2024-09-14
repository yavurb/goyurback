package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavurb/goyurback/internal/database/postgres"
	"github.com/yavurb/goyurback/internal/projects/domain"
)

type Repository struct {
	db *postgres.Queries
}

func NewRepo(connpool *pgxpool.Pool) domain.ProjectRepository {
	return &Repository{
		db: postgres.New(connpool),
	}
}

func (r *Repository) CreateProject(ctx context.Context, project *domain.ProjectCreate) (*domain.Project, error) {
	postID := pgtype.Int4{Valid: false}

	if _postID := project.PostID; _postID != 0 {
		postID = pgtype.Int4{Int32: project.PostID, Valid: true}
	}

	project_, err := r.db.CreateProject(ctx, postgres.CreateProjectParams{
		PublicID:     project.PublicID,
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

	newProject := toDomainStruct(&project_)

	return newProject, nil
}

func (r *Repository) GetProject(ctx context.Context, id string) (*domain.Project, error) {
	project_, err := r.db.GetProject(ctx, id)
	if err != nil {
		log.Printf("DB Error getting project: %v\n", err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProjectNotFound
		}

		return nil, err
	}

	project := toDomainStruct(&project_)

	return project, nil
}

func (r *Repository) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	projects_, err := r.db.GetProjects(ctx)
	if err != nil {
		log.Printf("DB Error obtaining projects: %v", err)

		if errors.Is(err, pgx.ErrNoRows) {
			return []*domain.Project{}, nil
		}

		return nil, err
	}

	projects := []*domain.Project{}

	for _, project := range projects_ {
		projects = append(projects, toDomainStruct(&project))
	}

	return projects, nil
}

func toDomainStruct(project_ *postgres.Project) *domain.Project {
	return &domain.Project{
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
}
