package application

import (
	"context"
	"log"

	"github.com/yavurb/goyurback/internal/projects/domain"
)

func (uc *projectUsecase) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	projects, err := uc.repository.GetProjects(ctx)

	if err != nil {
		log.Printf("Error getting projects. Got: %v\n", err)

		return nil, err
	}

	return projects, nil
}
