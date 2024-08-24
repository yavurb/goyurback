package application

import "github.com/yavurb/goyurback/internal/auth/domain"

type apiKeyUsecase struct {
	repository domain.APIKeyRepository
}

func NewAPIKeyUsecase(repository domain.APIKeyRepository) domain.APIKeyUsecase {
	return &apiKeyUsecase{repository}
}
