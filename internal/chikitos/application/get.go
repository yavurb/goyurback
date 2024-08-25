package application

import (
	"context"
	"log"

	"github.com/yavurb/goyurback/internal/chikitos/domain"
)

func (uc *ChikitoUsecase) Get(ctx context.Context, id string) (*domain.Chikito, error) {
	chikito, err := uc.repository.GetChikito(ctx, id)
	if err != nil {
		log.Printf("Unable to get chikito. Got: %v\n", err)

		return nil, domain.ErrChikitoNotFound
	}

	return chikito, err
}
