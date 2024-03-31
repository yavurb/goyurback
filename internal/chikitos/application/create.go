package application

import (
	"context"
	"log"

	"github.com/yavurb/goyurback/internal/chikitos/domain"
)

func (uc *ChikitoUsecase) Create(ctx context.Context, url, description string) (*domain.Chikito, error) {
	chikitoCreated, err := uc.repository.CreateChikito(ctx, url, description)
	if err != nil {
		log.Printf("Unable to create chikito. Got: %v\n", err)

		return nil, err
	}

	return chikitoCreated, nil
}
