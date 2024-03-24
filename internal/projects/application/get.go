package application

import (
	"context"
	"log"

	"github.com/yavurb/goyurback/internal/projects/domain"
)

func (uc *projectUsecase) Get(ctx context.Context, id string) (*domain.Project, error) {
	project, err := uc.repository.GetProject(ctx, id)

	if err != nil {
		log.Printf("Error getting project. Got: %v\n", err)

		return nil, domain.ErrProjectNotFound
	}

	return project, nil
}
