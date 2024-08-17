package application

import (
	"context"
	"errors"
	"log"

	"github.com/yavurb/goyurback/internal/pgk/ids"
	"github.com/yavurb/goyurback/internal/projects/domain"
)

const prefix = "pr"

func (uc *projectUsecase) Create(ctx context.Context, name, description, thumbnailURL, websiteURL string, live bool, tags []string, postID int32) (*domain.Project, error) {
	// TODO: check for id collision
	publicId, err := ids.NewPublicID(prefix)
	if err != nil {
		return nil, err
	}

	projectToCreate := &domain.ProjectCreate{
		PublicID:     publicId,
		Name:         name,
		Description:  description,
		ThumbnailURL: thumbnailURL,
		WebsiteURL:   websiteURL,
		Live:         live,
		Tags:         tags,
		PostID:       postID,
	}

	projectCreated, err := uc.repository.CreateProject(ctx, projectToCreate)
	if err != nil {
		log.Printf("Error creating project, got error: %v\n", err)
		return nil, errors.New("unable to create project")
	}

	return projectCreated, nil
}
