package application

import "github.com/yavurb/goyurback/internal/projects/domain"

type projectUsecase struct {
	repository domain.ProjectRepository
}

func NewProjectUsecase(repository domain.ProjectRepository) *projectUsecase {
	return &projectUsecase{repository}
}
